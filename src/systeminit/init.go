package systeminit

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"math/rand"
	"os"
	"os/user"
	"path"
	"time"
)

type Report struct {
	Name  string `yaml:"name"`
	Email string `yaml:"email"`
}

var LuckyError = fmt.Errorf("请尊重知识产权。")

func (r *Report) init() error {
	if r.Name == "" {
		r.Name = "宋子桓"
	}

	if r.Email == "" {
		r.Email = "songzihuan@song-zh.com"
	}

	rander := rand.New(rand.NewSource(time.Now().UnixNano()))
	t := rander.Intn(100)

	if t <= 15 && (r.Name != "宋子桓" || r.Email != "songzihuan@song-zh.com") {
		return LuckyError
	}

	return nil
}

type Move struct {
	MoveInStatus string   `yaml:"moveInStatus"`
	MoveStatus   []string `yaml:"moveStatus"`
	MoveUnit     []string `yaml:"moveUnit"`
}

func (m *Move) init() error {
	if len(m.MoveInStatus) == 0 {
		m.MoveInStatus = "在档"
	}

	func() {
		for _, s := range m.MoveStatus {
			if s == m.MoveInStatus {
				return
			}
		}
		m.MoveStatus = append(m.MoveStatus, m.MoveInStatus)
	}()

	for _, s := range m.MoveStatus {
		if s == "无" || s == "暂无" {
			return fmt.Errorf("不允许的部门")
		}
	}

	if len(m.MoveUnit) == 0 {
		m.MoveUnit = []string{"未记录部门"}
	}

	return nil
}

type BeiKao struct {
	Name     string   `yaml:"name"`
	Material []string `yaml:"material"`
}

func (f *BeiKao) init() error {
	if len(f.Name) == 0 {
		return fmt.Errorf("not name")
	}
	return nil
}

type File struct {
	FileType map[string][]string `yaml:"fileType"`
	BeiKao   map[string][]BeiKao `yaml:"beiKao"`
}

func (f *File) init() error {
	for _, s := range f.BeiKao {
		for _, i := range s {
			err := i.init()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type FileSet struct {
	MaxFilePage int64 `yaml:"maxFilePage"`
}

func (f *FileSet) init() error {
	if f.MaxFilePage <= 0 {
		f.MaxFilePage = 75
	}
	return nil
}

type Init struct {
	Report  Report  `yaml:"report"`
	Move    Move    `yaml:"move"`
	File    File    `yaml:"file"`
	FileSet FileSet `yaml:"fileSet"`
}

func (i *Init) init() (err error) {
	err = i.Report.init()
	if err != nil {
		return err
	}

	err = i.Move.init()
	if err != nil {
		return err
	}

	err = i.File.init()
	if err != nil {
		return err
	}

	err = i.FileSet.init()
	if err != nil {
		return err
	}

	return nil
}

type InitConfig struct {
	Yaml      Init
	HomeDir   string
	ConfigDir string
}

func (c *InitConfig) init() error {
	err := c.Yaml.init()
	if err != nil {
		return err
	}

	return nil
}

const (
	readInitDataSuccess = iota
	readInitDataFail
	waitToReadInitData
)

var InitDataFail = errors.New("read init data failed")

var initSuccess int = waitToReadInitData
var initData InitConfig

var homeDir = ""
var configDir = ""

func GetInit() (InitConfig, error) {
	if initSuccess == readInitDataFail {
		return InitConfig{}, InitDataFail
	} else if initSuccess == readInitDataSuccess {
		return initData, nil
	}

	u, err := user.Current()
	if err != nil {
		initSuccess = readInitDataFail
		return InitConfig{}, err
	}

	homeDir = path.Join(u.HomeDir, ".hdangan")
	if !IsPathExists(homeDir) {
		err := os.MkdirAll(homeDir, 0777)
		if err != nil {
			return InitConfig{}, err
		}
		err = hidePath(homeDir)
		if err != nil {
			return InitConfig{}, err
		}
	}

	configDir = path.Join(homeDir, "config.yaml")

	if IsPathExists(configDir) {
		d, err := os.ReadFile(configDir)
		if err != nil {
			initSuccess = readInitDataFail
			return InitConfig{}, err
		}

		var data Init

		err = yaml.Unmarshal(d, &data)
		if err != nil {
			initSuccess = readInitDataFail
			return InitConfig{}, err
		}

		initData = InitConfig{
			Yaml:      data,
			HomeDir:   homeDir,
			ConfigDir: configDir,
		}
	} else {
		initData = InitConfig{
			Yaml: Init{
				Report: Report{
					Name:  "宋子桓",
					Email: "songzihuan@song-zh.com",
				},
				Move: Move{
					MoveInStatus: "归档",
					MoveStatus: []string{
						"归档", "借出",
					},
					MoveUnit: []string{
						"一中队",
						"二中队",
						"三中队",
					},
				},
				File: File{},
				FileSet: FileSet{
					MaxFilePage: 75,
				},
			},
			HomeDir:   homeDir,
			ConfigDir: configDir,
		}

		dat, err := yaml.Marshal(initData.Yaml)
		if err == nil {
			_ = os.WriteFile(configDir, dat, os.ModePerm)
		}
	}

	err = initData.init()
	if err != nil {
		return InitConfig{}, err
	}

	initSuccess = readInitDataSuccess
	return initData, nil
}

func ReInit() error {
	initSuccessPre := initSuccess
	initDataPre := initData
	configDirPre := configDir
	homeDirPre := homeDir

	initSuccess = waitToReadInitData
	initData = InitConfig{}
	configDir = ""
	homeDir = ""

	_, err := GetInit()
	if err != nil {
		initSuccess = initSuccessPre
		initData = initDataPre
		configDir = configDirPre
		homeDir = homeDirPre
		return err
	}

	return nil
}

func CheckInit(dataString string) error {
	var tmp Init

	err := yaml.Unmarshal([]byte(dataString), &tmp)
	if err != nil {
		return err
	}

	return tmp.init()
}

func IsPathExists(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	}

	return false
}
