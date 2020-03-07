package machine

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/iesreza/foundation/lib"
	"github.com/iesreza/foundation/lib/network"
	"os/exec"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

type DiskDrive struct {
	Caption      string `json:"name"`
	DeviceID     string `json:"device_id"`
	Model        string `json:"model"`
	Partitions   uint   `json:"partitions"`
	Size         uint64 `json:"size"`
	SerialNumber string `json:"serial_number"`
	Active       bool   `json:"active"`
}

func UniqueHwID() (string, error) {
	netconfig, _ := network.GetConfig()
	mac := netconfig.HardwareAddress.String()
	hardDiskId, _ := GetActiveHddSerial()
	biosId, _ := GetBiosId()
	if len(mac) == 17 {
		mac = strings.Replace(mac, ":", "", -1)[6:]
	}
	if len(biosId) > 6 {
		biosId = biosId[len(biosId)-6:]
	}
	if len(hardDiskId) > 2 && hardDiskId[0:2] == "0x" {
		hardDiskId = hardDiskId[2:]
	}
	a := 0
	b := 0
	c := 0

	hwid := ""
	for {

		if a >= 0 && a < len(mac) {
			hwid += string(mac[a])
			a++
		} else {
			a = -1
		}
		if b >= 0 && b < len(hardDiskId) {
			hwid += string(hardDiskId[b])
			b++
		} else {
			b = -1
		}
		if c >= 0 && c < len(biosId) {
			hwid += string(biosId[c])
			c++
		} else {
			c = -1
		}
		if a == -1 && b == -1 && c == -1 {
			break
		}

	}

	reg := regexp.MustCompile("[^a-zA-Z0-9]+")
	hwid = reg.ReplaceAllString(hwid, "")

	return hwid, nil
}

func GetBiosId() (string, error) {

	if runtime.GOOS == "windows" {
		res, err := exec.Command(`cmd`, "/C", `wmic csproduct get UUID`).CombinedOutput()
		if err == nil {
			lines := strings.Split(strings.TrimSpace(string(res)), "\n")
			if len(lines) == 2 {
				return strings.TrimSpace(lines[1]), nil
			}
		}

	} else {
		res, err := exec.Command(`bash`, "-c", `dmidecode -s system-uuid`).CombinedOutput()

		if err == nil {
			return strings.TrimSpace(string(res)), nil
		}
	}
	return "", fmt.Errorf("unable to get bios id")

}

func GetActiveHddSerial() (string, error) {
	drive := GetActiveHardDisk()
	if drive.DeviceID == "" {
		return "", fmt.Errorf("active drive not found")

	}
	return drive.SerialNumber, nil
}

func GetHardDisks() ([]DiskDrive, error) {
	if runtime.GOOS == "windows" {
		res, err := exec.Command("bash", "-c", `wmic DiskDrive get Caption,DeviceID,Model,Partitions,Size,SerialNumber format:csv`).CombinedOutput()

		if err != nil {
			return []DiskDrive{}, err
		} else {
			dd := []DiskDrive{}
			err = parseWmicResult(string(res), &dd)
			fmt.Println(err)

			if err == nil {
				res, err := exec.Command(`cmd`, "/C", `wmic partition where BootPartition=true get DeviceID`).CombinedOutput()
				if err == nil {
					re := regexp.MustCompile(`(?m)Disk\s+\#(\d+)`)
					match := re.FindAllStringSubmatch(string(res), 1)
					if len(match) > 0 && len(match[0]) == 2 {
						re := regexp.MustCompile(`.+(\d+)$`)
						for k, item := range dd {
							matches := re.FindAllStringSubmatch(item.DeviceID, 1)
							if len(matches) == 1 && len(matches[0]) == 2 && matches[0][1] == match[0][1] {
								dd[k].Active = true
								break
							}
						}

					}
				}
			}

			return dd, err
		}
	} else {
		res, err := exec.Command("bash", "-c", "lsblk -o name,serial,size,mountpoint,vendor,model,type,kname,fstype -b -J").CombinedOutput()
		if err != nil {
			return []DiskDrive{}, err
		} else {
			dd := struct {
				Blockdevices []struct {
					Name       string `json:"name"`
					KernelName string `json:"kname"`
					Serial     string `json:"serial"`
					Size       string `json:"size"`
					Mountpoint string `json:"mountpoint"`
					Vendor     string `json:"vendor"`
					Model      string `json:"model"`
					Type       string `json:"type"`
					Fstype     string `json:"fstype"`
					Children   []struct {
						Name       string      `json:"name"`
						Fstype     interface{} `json:"fstype"`
						Mountpoint string      `json:"mountpoint"`
					} `json:"children,omitempty"`
				} `json:"blockdevices"`
			}{}

			json.Unmarshal(res, &dd)
			response := []DiskDrive{}

			for _, item := range dd.Blockdevices {

				if item.Type != "disk" || item.Mountpoint == "[SWAP]" {
					continue
				}

				drive := DiskDrive{}
				drive.Caption = item.Vendor + " " + item.Model
				drive.DeviceID = item.KernelName
				drive.Model = item.Vendor + " " + item.Model
				drive.SerialNumber = item.Serial
				drive.Size, _ = strconv.ParseUint(item.Size, 0, 64)
				drive.Partitions = uint(len(item.Children))
				for _, child := range item.Children {
					if child.Mountpoint == "/" {
						drive.Active = true
						break
					}
				}
				if drive.Size < 20*lib.MB {
					continue
				}
				response = append(response, drive)

			}
			return response, nil
		}
	}
	return []DiskDrive{}, nil
}

//Parse the csv format output of the RunCmd
func parseWmicResult(stdout string, dst interface{}) error {
	dv := reflect.ValueOf(dst).Elem()
	t := dv.Type().Elem()

	dv.Set(reflect.MakeSlice(dv.Type(), 0, 0))

	lines := strings.Split(stdout, "\n")
	var header []int = nil

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			v := reflect.New(t)
			r := csv.NewReader(strings.NewReader(line))

			r.FieldsPerRecord = t.NumField() + 1

			records, err := r.ReadAll()
			if err != nil {
				return err
			}
			//Find the field number of the record
			if header == nil {
				header = make([]int, len(records[0]), len(records[0]))
				for i, record := range records[0] {
					for j := 0; j < t.NumField(); j++ {
						if record == t.Field(j).Name {
							header[i] = j
						}
					}
				}
				continue
			} else {
				for i, record := range records[0] {
					f := reflect.Indirect(v).Field(header[i])
					switch t.Field(header[i]).Type.Kind() {
					case reflect.String:
						f.SetString(record)
					case reflect.Uint, reflect.Uint64:
						uintVal, err := strconv.ParseUint(record, 10, 64)
						if err != nil {
							return err
						}
						f.SetUint(uintVal)
					case reflect.Bool:
						bVal, err := strconv.ParseBool(record)
						if err != nil {
							return err
						}
						f.SetBool(bVal)
					default:
						return errors.New("unknown data type!")
					}
				}

			}
			dv.Set(reflect.Append(dv, reflect.Indirect(v)))
		}
	}

	return nil
}

func GetActiveHardDisk() DiskDrive {
	hdds, _ := GetHardDisks()

	for _, item := range hdds {
		if item.Active {
			return item
		}
	}

	return DiskDrive{}
}
