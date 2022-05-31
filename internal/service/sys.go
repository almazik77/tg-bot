package service

import (
	"bufio"
	"log"
	"market-bot/internal/service/model"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// Sys implements basic bot function to respond on ping and others from basic.data file.
// also, reacts on say! with keys/values from say.data file
type Sys struct {
	dataLocation string
	commands     []model.SysCommand
}

// NewSys makes new sys bot and load data to []say and basic map
func NewSys(dataLocation string) (*Sys, error) {
	log.Printf("[INFO] created sys bot, data location=%s", dataLocation)
	res := Sys{dataLocation: dataLocation}
	if err := res.loadBasicData(); err != nil {
		return nil, err
	}

	rand.Seed(0)
	return &res, nil
}

func (p Sys) OnMessage(text string) []model.SysCommand {
	var commands []model.SysCommand
	if text == "" {
		return p.commands
	}
	for _, bot := range p.commands {
		if found := containsLike(bot.Triggers, strings.ToLower(text)); found {
			commands = append(commands, bot)
		}
	}

	return commands
}

func (p *Sys) loadBasicData() error {
	bdata, err := readLines(p.dataLocation + "/basic.data")
	if err != nil {
		return errors.Wrap(err, "can't load basic.data")
	}

	for _, line := range bdata {
		elems := strings.Split(line, "|")
		if len(elems) != 3 {
			log.Printf("[DEBUG] bad format %s, ignored", line)
			continue
		}
		sysCommand := model.SysCommand{
			Question: elems[1],
			Message:  strings.ReplaceAll(elems[2], "\\n", "\n"),
			Triggers: strings.Split(elems[0], ";"),
		}
		p.commands = append(p.commands, sysCommand)
		log.Printf("[DEBUG] loaded basic response, %v, %s", sysCommand.Triggers, sysCommand.Message)
	}
	return nil
}

func readLines(path string) ([]string, error) {
	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, errors.Wrapf(err, "can't open %s", path)
	}
	defer f.Close() //nolint

	result := make([]string, 0)
	s := bufio.NewScanner(f)
	for s.Scan() {
		result = append(result, s.Text())
	}

	return result, nil
}

// ReactOn keys
func (p Sys) ReactOn() []string {
	res := make([]string, 0)
	for _, bot := range p.commands {
		res = append(bot.Triggers, res...)
	}
	return res
}

func containsLike(s []string, e string) bool {
	e = strings.TrimSpace(e)
	for _, a := range s {
		if strings.Contains(a, e) {
			return true
		}
	}
	return false
}
