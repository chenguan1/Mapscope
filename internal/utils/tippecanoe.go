package utils

import (
	"Mapscope/internal/models"
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"os/exec"
)

func CreateMbtiles(gjsons []string, name string, mbfile string, task *models.Task) error  {
	var params []string
	params = append(params, "-zg")
	params = append(params, "-o")
	params = append(params, mbfile)
	params = append(params, "--force")
	params = append(params, []string{"-n", name}...)
	params = append(params, []string{"-l", name}...)
	params = append(params, "-P")
	params = append(params, "--drop-densest-as-needed")
	params = append(params, "--extend-zooms-if-still-dropping")
	params = append(params, []string{"-t", "./"}...)
	params = append(params, gjsons...)

	//logrus.Info(params)
	if err := os.Remove(mbfile); err != nil{
		logrus.Errorf("remove file error, %v",err)
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd := exec.Command("d:\\tippecanoe\\tippecanoe", params...)
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()
	// var errStdout, errStderr error
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("CreateMbtiles failed: %v", err)
	}
	go func() {
		io.Copy(stdout, stdoutIn)
	}()
	go func() {
		io.Copy(stderr, stderrIn)
	}()

	task.Progress = 50

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("CreateMbtiles process failed: %v", err)
	}

	return nil
}