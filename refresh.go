package main

import (
	"os/exec"
	"log"
	"bytes"
        "strings"
	"os"
	"io/ioutil"
)

func Refresh() {

        // get ipv4 pool of cluster through calico
        var ipPool []string
	ipPool = getIpPool()

	// get ip pool config file content 
	content, _ := getIpPoolConfig("/etc/cni/net.d/10-calico.conf", ipPool)

	// create a new config file in direct /etc/cni/net.d/new.conf
        var configContent string
        configContent = strings.Join(content, "")
	f, err := os.Create("/etc/cni/net.d/new.conf")
        check(err)

        defer f.Close()

        n, err := f.WriteString(configContent)
        check(err)
        log.Printf("finish updating config file %d bytes \n", n)

	f.Sync()

}

func getIpPool() []string {

	path, err := exec.LookPath("calicoctl")
	if err != nil {
		log.Fatal("calicoctl not found")
	}
	log.Printf(`"calicoctl" is available at %s `, path)

	cmd := exec.Command("calicoctl", "get", "ipPool")
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Printf("Unprocess output is: \n %q", out.String())

        var tmpIpPool []string
        var s string
        s = out.String()

        for _, lineStr := range strings.Split(s, "\n") {
		lineStr = strings.TrimSpace(lineStr)
		if (lineStr == "" || strings.Contains(lineStr, "CIDR") || strings.Contains(lineStr, "fd80")) {
			continue
		}else {
			tmpIpPool = append(tmpIpPool, lineStr)
		}
	}

	log.Printf("read latest ipv4 pool of cluster successfully.")

        return tmpIpPool
}

func getIpPoolConfig(filePath string, info []string) ([]string, error) {

	result := []string{}
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("read file: %v error: %v", filePath, err)
		return result, err
	}
	s := string(b)
        log.Printf("upload config file successfully: %v, size %v", filePath, len(s))

	// fmt.Println("Source content:", s)
	for _, lineStr := range strings.Split(s, "\n") {
		//lineStr = strings.TrimSpace(lineStr)
		if lineStr == "" {
		continue
		}

		if strings.Contains(lineStr, "ipv4_pools") {
			tmpInfo := strings.Join(info, `","`)
			tmpLineStr := make([]string, 1)
			tmpLineStr[0] = `	"ipv4_pools": ["`
			tmpLineStr = append(tmpLineStr, tmpInfo)
			tmpLineStr = append(tmpLineStr, `"]`)
			lineStr = strings.Join(tmpLineStr, "")
		}

		result = append(result, lineStr)
		result = append(result, "\n")
	}

	log.Printf("get ip pool config file content successfully: %v, size %v", filePath, len(result))

	return result, nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
