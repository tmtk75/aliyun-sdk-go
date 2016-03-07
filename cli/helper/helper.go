package helper

import (
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/tmtk75/cli"
)

func Atoi(a string) int {
	i, err := strconv.Atoi(a)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	return i
}

func Merge(a, b []cli.Command) []cli.Command {
	mv := func(a []cli.Command) map[string]cli.Command {
		p := make(map[string]cli.Command)
		for _, e := range a {
			p[e.Name] = e
		}
		return p
	}

	mg := func(a, b cli.Command) cli.Command {
		if b.Usage != "" {
			a.Usage = b.Usage
		}
		if b.Action != nil {
			a.Action = b.Action
		}
		if b.Args != "" {
			a.Args = b.Args
		}
		if b.Flags != nil {
			a.Flags = b.Flags
		}
		return a
	}

	p := mv(a)
	q := mv(b)
	d := make(cmds, 0)
	for k, c := range p {
		if v, ok := q[k]; ok {
			d = append(d, mg(c, v))
		} else {
			d = append(d, c)
		}
	}

	sort.Sort(d)
	return d
}

type cmds []cli.Command

func (c cmds) Len() int {
	return len(c)
}

func (c cmds) Less(i, j int) bool {
	return strings.Compare(c[i].Name, c[j].Name) < 1
}

func (c cmds) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func parseTime(s string) time.Time {
	if s == "" {
		return time.Time{}
	}

	//t, err := time.Parse("2006-01-02T15:04:05Z07:00", s)
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	return t
}

func EndTime(s string) time.Time {
	t := parseTime(s)
	//fmt.Printf("given str: %v\n", s)
	//fmt.Printf("parsed: %v\n", t)
	return t
}

// Make StartTime with command line options
func StartTime(c *cli.Context) time.Time {
	st := c.GlobalString("start-time")
	if st != "" {
		return parseTime(st)
	}

	etstr := c.GlobalString("end-time")
	//fmt.Printf("given end-time: %v\n", etstr)
	et := parseTime(etstr)
	//fmt.Printf("end-time: %v\n", et.Format("2006-01-02T15:04:03Z"))
	//fmt.Printf("start-time: %v\n", st.Format("2006-01-02T15:04:03Z"))
	if v := c.GlobalString("duration"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			log.Fatalf("%v\n", err)
		}
		return et.Add(-d)
	}

	return time.Time{}
}
