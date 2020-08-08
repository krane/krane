package spec

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/biensupernice/krane/internal/storage"
	"github.com/biensupernice/krane/internal/utils"
)

var (
	Collection = "specs"
)

type Spec struct {
	Name      string `json:"name" binding:"required"`
	Config    Config `json:"config" binding:"required"`
	UpdateAt  string `json:"updated_at"`
	CreatedAt string `json:"created_at"`
}

type Config struct {
	Registry      string            `json:"registry"`
	Image         string            `json:"image" binding:"required"`
	ContainerPort string            `json:"container_port"`
	HostPort      string            `json:"host_port"`
	Env           map[string]string `json:"env"`
	Tag           string            `json:"tag"`
	Volumes       map[string]string `json:"volumes"`
}

func (s *Spec) CreateSpec() error {
	err := s.Validate()
	if err != nil {
		return err
	}

	s.setDefaults()
	s.CreatedAt = utils.UTCDateString()
	s.UpdateAt = utils.UTCDateString()

	data, err := storage.Get(Collection, s.Name)
	if err != nil {
		return err
	}

	if data != nil {
		return errors.New("spec already exists")
	}

	bytes, err := json.Marshal(s)
	if err != nil {
		return err
	}

	err = storage.Put(Collection, s.Name, bytes)
	if err != nil {
		return err
	}

	return nil
}

func (s *Spec) UpdateSpec(prevDeploymentName string) error {
	err := s.Validate()
	if err != nil {
		return err
	}

	prevSpec, err := GetOne(prevDeploymentName)
	if err != nil {
		return err
	}


	s.setDefaults()
	s.UpdateAt = utils.UTCDateString()
	s.CreatedAt = prevSpec.CreatedAt;

	bytes, err := json.Marshal(s)
	if err != nil {
		return err
	}

	err = storage.Put(Collection, s.Name, bytes)
	if err != nil {
		return err
	}

	// Only remove the old spec if the names are variant. This is because the name is what is used as the key
	if s.Name != prevDeploymentName {
		err = storage.Remove(Collection, prevDeploymentName)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetOne(name string) (Spec, error) {
	bytes, err := storage.Get(Collection, name)
	if err != nil {
		return Spec{}, err
	}

	if len(bytes) == 0 {
		return Spec{}, fmt.Errorf("Spec not found")
	}

	// Unmarshal bytes into Spec
	var s Spec
	json.Unmarshal(bytes, &s)

	return s, nil
}

func GetAll() ([]Spec, error) {
	bytes, err := storage.GetAll(Collection)
	if err != nil {
		return make([]Spec, 0), err
	}

	specs := make([]Spec, 0)
	for _, spec := range bytes {
		var s Spec
		err := json.Unmarshal(spec, &s)
		if err != nil {
			return specs, err
		}

		specs = append(specs, s)
	}

	return specs, nil
}

func (s Spec) Delete() error { return storage.Remove(Collection, s.Name) }

func (s Spec) Validate() error {
	isValidSpecName := s.isValidSpecName()
	if !isValidSpecName {
		return errors.New("Invalid Spec name")
	}

	return nil
}

func (s Spec) isValidSpecName() bool {
	startsWithLetter := "[a-z]"
	allowedCharacters := "[a-z0-9_-]"
	endWithLowerCaseAlphanumeric := "[0-9a-z]"
	characterLimit := "{2,}"

	matchers := fmt.Sprintf(`^%s%s*%s%s$`, // ^[a - z][a - z0 - 9_ -]*[0-9a-z]$
		startsWithLetter,
		allowedCharacters,
		endWithLowerCaseAlphanumeric,
		characterLimit)

	match := regexp.MustCompile(matchers)
	return match.MatchString(s.Name)
}

func (s *Spec) setDefaults() {
	if s.Config.Registry == "" {
		s.Config.Registry = "docker.io"
	}

	if s.Config.Env == nil {
		s.Config.Env = make(map[string]string, 0)
	}

	if s.Config.Volumes == nil {
		s.Config.Volumes = make(map[string]string, 0)
	}

	if s.Config.Tag == "" {
		s.Config.Tag = "latest"
	}
}
