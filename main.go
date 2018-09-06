package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var countsuccess int

func GetFilesAndDirsEx(dirPth string, files *[]string, dirs *[]string) (err error) {
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return err
	}

	PthSep := string(os.PathSeparator)

	for _, fi := range dir {
		if fi.IsDir() { // 目录, 递归遍历
			*dirs = append(*dirs, dirPth+PthSep+fi.Name())
			GetFilesAndDirsEx(dirPth+PthSep+fi.Name(), files, dirs)
		} else {
			// 过滤指定格式
			ok := strings.HasSuffix(fi.Name(), ".zip")
			if ok {
				if strings.Contains(fi.Name(), "-target_files-") && strings.Contains(dirPth, "official_") {
					*files = append(*files, dirPth+PthSep+fi.Name())
				}
			}
		}
	}

	return nil
}

func GetSetupDir(filepath string) string {
	PthSep := string(os.PathSeparator)
	split_strings := strings.Split(filepath, PthSep)
	len1 := len(split_strings)
	if len1 >= 3 {
		return split_strings[len1-3]
	}
	fmt.Printf("error, filepath is invalid,file is :[%s]", filepath)
	return ""
}

func exe_cmd(command string, wg *sync.WaitGroup) {
	fmt.Printf("Execute Shell start:[%s]!!!!!!\n", command)
	cmd := exec.Command("/bin/bash", "-c", command)
	_, err := cmd.Output()
	if err != nil {
		fmt.Printf("Execute Shell:%s failed with error:%s", command, err.Error())
		return
	}
	fmt.Printf("Execute Shell success!!!!!!\n")
	wg.Done()
}

func main() {
	var files []string = make([]string, 0)
	var dirs []string = make([]string, 0)
	GetFilesAndDirsEx("/data/agoldbase_rom", &files, &dirs)

	var maxfile string
	var dst_file string
	for _, file := range files {
		file_day := GetSetupDir(file)
		if file_day > maxfile {
			maxfile = file_day
			dst_file = file
		}
		//fmt.Printf("获取的文件为[%s]\n", file)
	}
	if dst_file == "" {
		fmt.Printf("基准文件文件获取失败\n")
		return
	}
	fmt.Printf("基准文件为[%s]\n", dst_file)

	PthSep := string(os.PathSeparator)
	string_temp1 := strings.Split(dst_file, PthSep)
	len1 := len(string_temp1)
	name1 := string_temp1[len1-3]

	wg := new(sync.WaitGroup)
	countall := 0
	countsuccess = 0
	for _, file := range files {
		if file == dst_file {
			continue
		}
		string_temp2 := strings.Split(file, PthSep)
		len2 := len(string_temp2)
		name2 := string_temp2[len2-3]
		updatefile_name := "rom_" + name2 + "_" + name1 + ".update.zip"
		updatefile_name = "/data/wangcanli/official_update.zip/" + updatefile_name

		commands := "/data/official/build/tools/releasetools/ota_from_target_files --block -k "
		commands += "/data/official/build/target/product/security/releasekey  -s "
		commands += "/data/official/device/mediatek/build/releasetools/mt_ota_from_target_files.py -v --revision_info "
		commands += "/data/wangcanli/OTA/info.txt -i "
		commands += file
		commands += " " + dst_file
		commands += " " + updatefile_name

		wg.Add(1)
		countall++
		go exe_cmd(commands, wg)
		//break
	}
	fmt.Printf("all update file num is %d\n", countall)
	wg.Wait()
	fmt.Printf("buid update file end[%d:%d]\n", countsuccess, countall)
}
