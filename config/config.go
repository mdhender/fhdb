/*
 * Copyright (c) 2021 Michael D Henderson
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package config

import (
	"flag"
	"fmt"
	"github.com/peterbourgon/ff/v3"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	Serve  bool
	Debug  bool
	Data   string
	Server struct {
		Scheme  string
		Host    string
		Port    int
		Timeout struct {
			Idle  time.Duration
			Read  time.Duration
			Write time.Duration
		}
		JWT struct {
			Key string
		}
		TLS struct {
			Serve    bool
			CertFile string
			KeyFile  string
		}
		Web struct {
			Root string
		}
	}
}

// Default returns a default configuration.
// These are the values without loading the environment, configuration file, or command line.
func Default() *Config {
	var cfg Config
	cfg.Data = "D:\\GoLand\\fhdb\\data"
	cfg.Server.Scheme = "http"
	cfg.Server.Host = "localhost"
	cfg.Server.Port = 10801
	cfg.Server.Timeout.Idle = 10 * time.Second
	cfg.Server.Timeout.Read = 5 * time.Second
	cfg.Server.Timeout.Write = 10 * time.Second
	cfg.Server.Web.Root = "D:\\GoLand\\fhdb\\web"
	return &cfg
}

// Load updates the values in a Config in this order:
//   1. It will load a configuration file if one is given on the
//      command line via the `-config` flag. If provided, the file
//      must contain a valid JSON object.
//   2. Environment variables, using the prefix `CONDUIT_RYER_SERVER`
//   3. Command line flags
func (cfg *Config) Load() error {
	fs := flag.NewFlagSet("Server", flag.ExitOnError)
	debug := fs.Bool("debug", cfg.Debug, "log debug information (optional)")
	data := fs.String("data", cfg.Data, "path to application data")
	serverScheme := fs.String("scheme", cfg.Server.Scheme, "http scheme, either 'http' or 'https'")
	serverHost := fs.String("host", cfg.Server.Host, "host name (or IP) to listen on")
	serverPort := fs.Int("port", cfg.Server.Port, "port to listen on")
	serverJWTKey := fs.String("jwt-key", cfg.Server.JWT.Key, "jwt hs256 key")
	serverTimeoutIdle := fs.Duration("idle-timeout", cfg.Server.Timeout.Idle, "http idle timeout")
	serverTimeoutRead := fs.Duration("read-timeout", cfg.Server.Timeout.Read, "http read timeout")
	serverTimeoutWrite := fs.Duration("write-timeout", cfg.Server.Timeout.Write, "http write timeout")
	serverTLSServe := fs.Bool("https", cfg.Server.TLS.Serve, "serve https")
	serverTLSCertFile := fs.String("https-cert-file", cfg.Server.Host, "https certificate file")
	serverTLSKeyFile := fs.String("https-key-file", cfg.Server.Host, "https certificate key file")
	serverWebRoot := fs.String("web", cfg.Server.Web.Root, "path to serve assets from")

	if err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("FHOE"), ff.WithConfigFileFlag("config"), ff.WithConfigFileParser(ff.JSONParser)); err != nil {
		return err
	}

	cfg.Debug = *debug
	cfg.Data = filepath.Clean(*data)
	cfg.Server.Scheme = *serverScheme
	cfg.Server.Host = *serverHost
	cfg.Server.Port = *serverPort
	cfg.Server.JWT.Key = *serverJWTKey
	cfg.Server.Timeout.Idle = *serverTimeoutIdle
	cfg.Server.Timeout.Read = *serverTimeoutRead
	cfg.Server.Timeout.Write = *serverTimeoutWrite
	cfg.Server.TLS.Serve = *serverTLSServe
	cfg.Server.TLS.CertFile = *serverTLSCertFile
	cfg.Server.TLS.KeyFile = *serverTLSKeyFile
	cfg.Server.Web.Root = *serverWebRoot

	if cfg.Server.TLS.Serve == true {
		if cfg.Server.TLS.CertFile == "" {
			return fmt.Errorf("must supply certificates file when serving HTTPS")
		}
		if cfg.Server.TLS.KeyFile == "" {
			return fmt.Errorf("must supply certificate key file when serving HTTPS")
		}
	}

	return nil
}
