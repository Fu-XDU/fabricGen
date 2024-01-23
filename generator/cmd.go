package generator

import (
	"log"
	"os/exec"
)

// ExecShell 阻塞式的执行外部shell命令的函数,等待执行完毕并返回标准输出
func ExecShell(s string) (string, error) {
	// 函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command("/bin/bash", "-c", s)
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Fatalf("Exec shell cmd `%s` failed, output:%v, error: %v", cmd, string(output), err)
	}
	return string(output), err
}
