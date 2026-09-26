// Конфигуратор приложения с `sync.Once`
package main

cm.config = map[string]string{
"app_name":  "MyApp",
"port":      "8080",
"log_level": "debug",
}

func main() {
    keys := []string{"app_name", "port", "log_level"}
    configManager.PrintConfig()
}
