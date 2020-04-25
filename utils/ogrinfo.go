package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type OgrinfoPgParams struct {
	Host string
	Port string
	DbName string
	Username string
	Password string
}

type OgrInfo struct {
}


func OgrinfoPg(p OgrinfoPgParams, layerName string) (*OgrInfo, error) {
	//ogrinfo -so "PG:dbname=mapscope host=localhost port=5432 user=postgres password=111111" dataset_test0

	var pms []string
	pg := fmt.Sprintf(`PG:dbname=%s host=%s port=%s user=%s password0=%s`,
		p.DbName, p.Host, p.Port, p.Username, p.Password)
	pms = append(pms, "-so")
	pms = append(pms, pg)
	pms = append(pms, layerName)

	cmd := exec.Command("ogrinfo", pms...)

	var stdoutBuf, stderrBuf bytes.Buffer
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)
	go func() {
		io.Copy(stdout, stdoutIn)
	}()
	go func() {
		io.Copy(stderr, stderrIn)
	}()

	err := cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("get dataset info failed. err: %v", err)
	}

	err = cmd.Wait()
	if err != nil {
		return nil, fmt.Errorf("wait dataset info failed. err: %v", err)
	}

	rawinfo := stdoutBuf.String()
	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++")
	fmt.Println(rawinfo)


	return nil, nil
}