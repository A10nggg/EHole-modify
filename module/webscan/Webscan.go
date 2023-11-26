package webScan

import (
	"embed"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"ehole/module/webscan/common"
	"ehole/module/webscan/lib"
)

//go:embed pocs
var Pocs embed.FS
var once sync.Once
var AllPocs []*lib.Poc

func WebScan(targets []string) {
	//lib.Inithttp(common.Pocinfo)
	//fmt.Printf("-%v", info)
	once.Do(initpoc)
	//var pocinfo = common.Pocinfo
	// buf := strings.Split(info.Url, "/")
	// pocinfo.Target = strings.Join(buf[:3], "/")

	for _, target := range targets {
		//pocinfo.PocName = lib.CheckInfoPoc(infostr)
		Execute(target)
	}

}

func Execute(target string) {
	req, err := http.NewRequest("GET", target, nil)
	if err != nil {
		errlog := fmt.Sprintf("[-] webpocinit %v %v", target, err)
		common.LogError(errlog)
		return
	}
	req.Header.Set("User-agent", common.UserAgent)
	req.Header.Set("Accept", common.Accept)
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")

	req.Header.Set("Connection", "close")
	pocs := filterPoc("")

	lib.CheckMultiPoc(req, pocs, 20)
}

func initpoc() {
	if common.PocPath == "" {
		entries, err := Pocs.ReadDir("pocs")
		if err != nil {
			fmt.Printf("[-] init poc error: %v", err)
			return
		}
		for _, one := range entries {
			path := one.Name()
			if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
				if poc, _ := lib.LoadPoc(path, Pocs); poc != nil {
					AllPocs = append(AllPocs, poc)
				}
			}
		}
	} else {
		err := filepath.Walk(common.PocPath,
			func(path string, info os.FileInfo, err error) error {
				if err != nil || info == nil {
					return err
				}
				if !info.IsDir() {
					if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
						poc, _ := lib.LoadPocbyPath(path)
						if poc != nil {
							AllPocs = append(AllPocs, poc)
						}
					}
				}
				return nil
			})
		if err != nil {
			fmt.Printf("[-] init poc error: %v", err)
		}
	}
}

func filterPoc(pocname string) (pocs []*lib.Poc) {
	if pocname == "" {
		return AllPocs
	}
	for _, poc := range AllPocs {
		if strings.Contains(poc.Name, pocname) {
			pocs = append(pocs, poc)
		}
	}
	return
}
