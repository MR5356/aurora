package terminal

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"os"
	"path/filepath"
)

type KubernetesTerminal struct {
	config     Config
	kc         KubernetesConfig
	clientSet  *kubernetes.Clientset
	restConfig *rest.Config
	tty        remotecommand.Executor
	stdin      io.WriteCloser
	stdout     io.Reader
	sizeQueue  *SizeQueue
}

type KubernetesConfig struct {
	Namespace     string
	PodName       string
	ContainerName string
	Command       []string
}

type SizeQueue struct {
	resizeChan chan *remotecommand.TerminalSize
}

func NewSizeQueue() *SizeQueue {
	return &SizeQueue{
		resizeChan: make(chan *remotecommand.TerminalSize, 1),
	}
}

func (q *SizeQueue) Next() *remotecommand.TerminalSize {
	return <-q.resizeChan
}

func NewKubernetesTerminal(kc KubernetesConfig, config Config) *KubernetesTerminal {
	return &KubernetesTerminal{
		config:    config,
		kc:        kc,
		sizeQueue: NewSizeQueue(),
	}
}

func (t *KubernetesTerminal) Start() error {
	config, err := rest.InClusterConfig()
	if err != nil {
		// TODO 尝试从默认位置获取，后期需要改
		config, err = clientcmd.BuildConfigFromFlags("", filepath.Join(os.Getenv("HOME"), ".kube", "pi5-config"))
		if err != nil {
			return err
		}
	}

	t.restConfig = config
	config.TLSClientConfig.Insecure = true
	config.CAData = nil
	config.CAFile = ""

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	t.clientSet = clientset

	if ns, err := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{}); err == nil {
		for _, n := range ns.Items {
			logrus.Infof("namespace: %s", n.Name)
		}
	}

	if err := t.verify(); err != nil {
		return err
	}

	req := t.clientSet.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(t.kc.PodName).
		Namespace(t.kc.Namespace).
		SubResource("exec")

	req.VersionedParams(&corev1.PodExecOptions{
		Container: t.kc.ContainerName,
		Command:   t.kc.Command,
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
		TTY:       true,
	}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(t.restConfig, "POST", req.URL())
	if err != nil {
		return err
	}

	t.tty = exec

	// 创建管道用于stdin和stdout
	stdinR, stdinW := io.Pipe()
	stdoutR, stdoutW := io.Pipe()

	go func() {
		t.tty.Stream(remotecommand.StreamOptions{
			Stdin:             stdinR,
			Stdout:            stdoutW,
			Stderr:            stdoutW,
			Tty:               true,
			TerminalSizeQueue: t.sizeQueue,
		})
	}()

	t.stdin = stdinW
	t.stdout = stdoutR

	return t.Resize(t.config.Cols, t.config.Rows)
}

func (t *KubernetesTerminal) Close() error {
	if t.stdin != nil {
		t.stdin.Close()
	}
	return nil
}

func (t *KubernetesTerminal) Write(p []byte) (n int, err error) {
	return t.stdin.Write(p)
}

func (t *KubernetesTerminal) Read(p []byte) (n int, err error) {
	return t.stdout.Read(p)
}

func (t *KubernetesTerminal) Resize(cols, rows uint32) error {
	size := &remotecommand.TerminalSize{
		Width:  uint16(cols),
		Height: uint16(rows),
	}

	select {
	case t.sizeQueue.resizeChan <- size:
		return nil
	default:
		select {
		case <-t.sizeQueue.resizeChan:
		default:
		}
		t.sizeQueue.resizeChan <- size
		return nil
	}
}

func (t *KubernetesTerminal) verify() error {
	// 获取 Pod
	pod, err := t.clientSet.CoreV1().Pods(t.kc.Namespace).Get(context.Background(), t.kc.PodName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("Pod %s/%s 不存在: %v", t.kc.Namespace, t.kc.PodName, err)
	}

	// 检查 Pod 状态
	if pod.Status.Phase != corev1.PodRunning {
		return fmt.Errorf("Pod %s/%s 不在运行状态，当前状态: %s", t.kc.Namespace, t.kc.PodName, pod.Status.Phase)
	}

	// 检查容器是否存在
	containerExists := false
	for _, container := range pod.Spec.Containers {
		if container.Name == t.kc.ContainerName {
			containerExists = true
			break
		}
	}
	if !containerExists {
		return fmt.Errorf("容器 %s 在 Pod %s/%s 中不存在", t.kc.ContainerName, t.kc.Namespace, t.kc.PodName)
	}
	return nil
}
