package revealgo

import (
	"fmt"
	"net"
	"os"
	"reflect"
	"strings"

	flags "github.com/jessevdk/go-flags"
)

type CLI struct {
}

type CLIOptions struct {
	Port       int    `short:"p" long:"port" description:"tcp port number of this server. default is 3000."`
	Theme      string `long:"theme" description:"slide theme or original css file name. default themes: beige, black, blood, league, moon, night, serif, simple, sky, solarized, and white" default:"black.css"`
	Transition string `long:"transition" description:"transition effect for slides: default, cube, page, concave, zoom, linear, fade, none" default:"default"`
	Watch      bool   `long:"watch" description:"watch for md file change and reload" default:"false"`
}

func (cli *CLI) Run() {
	opts, args, err := parseOptions()
	if err != nil {
		fmt.Printf("error:%v\n", err)
		os.Exit(1)
	}
	if len(args) < 1 {
		showHelp()
		os.Exit(0)
	}

	server := Server{
		port: opts.Port,
	}
	_, err = os.Stat(opts.Theme)
	originalTheme := false
	if err == nil {
		originalTheme = true
	}

	param := ServerParam{
		Path:          args[0],
		Theme:         addExtention(opts.Theme, "css"),
		OriginalTheme: originalTheme,
		Host:          GetLocalIP(),
		Watch:         opts.Watch,
		RevealOptions: map[string]interface{}{
			"transition": "'" + opts.Transition + "'",
		},
	}

	server.Serve(param)
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
func showHelp() {
	fmt.Fprint(os.Stderr, `Usage: revealgo [options] [MARKDOWN FILE]

Options:
`)
	t := reflect.TypeOf(CLIOptions{})
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag
		var o string
		if s := tag.Get("short"); s != "" {
			o = fmt.Sprintf("-%s, --%s", tag.Get("short"), tag.Get("long"))
		} else {
			o = fmt.Sprintf("--%s", tag.Get("long"))
		}
		fmt.Fprintf(os.Stderr, "  %-21s %s\n", o, tag.Get("description"))
	}
}

func parseOptions() (*CLIOptions, []string, error) {
	opts := &CLIOptions{}
	p := flags.NewParser(opts, flags.PrintErrors)
	args, err := p.Parse()
	if err != nil {
		return nil, nil, err
	}
	return opts, args, nil
}

func addExtention(path string, ext string) string {
	if strings.HasSuffix(path, fmt.Sprintf(".%s", ext)) {
		return path
	}
	path = fmt.Sprintf("%s.%s", path, ext)
	return path
}
