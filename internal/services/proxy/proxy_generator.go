package proxy

import (
	"fmt"
	"os"
	"path/filepath"
)

// GenerateNginxProxyConfig generates an Nginx proxy configuration for the given IP address, port, and subdomain,
// and writes it to the specified folder.
func GenerateNginxProxyConfig(ip string, port int, subdomain string) error {
	config := fmt.Sprintf(`
		server {
			listen 80;
			server_name %s;

			location / {
				proxy_pass http://%s:%d;
				proxy_set_header Host $http_host;
				proxy_set_header Upgrade $http_upgrade;
				proxy_set_header Connection upgrade;
				proxy_set_header Accept-Encoding gzip;
			}
		
		}
`, subdomain, ip, port)
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	folder := fmt.Sprintf("%s_proxy", ex)
	file_name := fmt.Sprintf("%s.conf", subdomain)
	filePath := filepath.Join(folder, file_name)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(config)
	if err != nil {
		return err
	}

	return nil
}
