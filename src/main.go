package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/rtitz/mac-ramdisk-creator/variables"
)

func hdiUtilAttach(size int) (string, error) {
	var output string
	cmd := exec.Command(variables.HdiUtil, "attach", "-nomount", "ram://"+strconv.Itoa(size))
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return output, errors.New("FAILED TO EXECUTE HDIUTIL ATTACH")
	}
	/*item := bufio.NewScanner(strings.NewReader(out.String()))
	for item.Scan() {
		fmt.Println(item.Text())
	}*/
	output = out.String()
	output = strings.TrimSpace(output)
	return output, nil
}

func hdiUtilDetach(device string) error {
	cmd := exec.Command(variables.HdiUtil, "detach", device)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return errors.New("FAILED TO EXECUTE HDIUTIL DETACH")
	}
	/*item := bufio.NewScanner(strings.NewReader(out.String()))
	for item.Scan() {
		fmt.Println(item.Text())
	}*/
	return nil
}

func diskInfo(name string) (string, error) {
	var output string
	cmd := exec.Command(variables.DiskUtil, "info", name)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return output, errors.New("DISK NOT FOUND")
	}
	item := bufio.NewScanner(strings.NewReader(out.String()))
	for item.Scan() {
		if strings.Contains(item.Text()+"\n", "APFS Physical Store:") {
			output = item.Text()
		}
	}
	output = strings.TrimSpace(output)
	var re1 = regexp.MustCompile(`.*disk`)
	var re2 = regexp.MustCompile(`s[1-9]+`)
	output = re1.ReplaceAllString(output, "/dev/disk")
	output = re2.ReplaceAllString(output, "")
	return output, nil
}

func ejectRamdisk(device string) error {
	cmd := exec.Command(variables.DiskUtil, "eject", device)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return errors.New("FAILED TO EJECT RAMDISK")
	}
	/*item := bufio.NewScanner(strings.NewReader(out.String()))
	for item.Scan() {
		fmt.Println(item.Text())
	}*/
	return nil
}

func createRamdisk(size int, name, hdiUtil_output string) error {
	//fmt.Printf("\nhdiutil_out: '%s'\nname: '%s'\n", hdiUtil_output, name)
	cmd := exec.Command(variables.DiskUtil, "partitionDisk", hdiUtil_output, "1", "GPTFormat", "APFS", name, "100%")

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return errors.New("FAILED TO CREATE RAMDISK")
	}
	/*item := bufio.NewScanner(strings.NewReader(out.String()))
	for item.Scan() {
		fmt.Println(item.Text())
	}*/
	return nil
}

func calculateSizeForDiskutil(sizeInMB int) int {
	return sizeInMB * 2000
}

func main() {
	// Define and check parameters
	sizeInMB := flag.Int("size", 0, "Specify the size of the ramdisk in MB. (Minimum is 10 MB)")
	name := flag.String("name", "", "Specify name of the ramdisk")
	eject := flag.Bool("eject", false, "Set this parameter to eject a previous created ramdisk")
	flag.Parse()

	if (*sizeInMB == 0 || *sizeInMB < 10) && !*eject {
		fmt.Printf("Parameter missing / wrong! Try again and specify the following parameters.\n\nParameter list:\n\n")
		flag.PrintDefaults()
		fmt.Printf("\n")
		os.Exit(999)
	}
	if *name == "" {
		*name = "ramdisk"
	}
	// End of: Define and check parameters

	fmt.Printf("%s %s\n", variables.AppName, variables.AppVersion)

	if !*eject {
		fmt.Printf("Creating ramdisk '%s' with %s MB ... ", *name, strconv.Itoa(*sizeInMB))
	} else {
		fmt.Printf("Ejecting ramdisk '%s' ... ", *name)
	}

	if !*eject {
		size := calculateSizeForDiskutil(*sizeInMB)

		hdiUtil_output, errhdiUtil := hdiUtilAttach(size)
		if errhdiUtil == nil {
			errCreateRamdisk := createRamdisk(size, *name, hdiUtil_output)
			if errCreateRamdisk != nil {
				_ = hdiUtilDetach(hdiUtil_output)
				log.Fatal(errCreateRamdisk.Error())
			}
		} else {
			_ = hdiUtilDetach(hdiUtil_output)
			log.Fatal(errhdiUtil.Error())
		}
	} else {
		device, errDiskInfo := diskInfo(*name)
		if errDiskInfo != nil {
			log.Fatal(errDiskInfo.Error())
		}
		errEjectRamdisk := ejectRamdisk(device)
		if errEjectRamdisk != nil {
			log.Fatal(errEjectRamdisk.Error())
		}
	}
	fmt.Printf("Finished!\n")
	os.Exit(0)
}
