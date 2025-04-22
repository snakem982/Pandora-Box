package systray

import (
	"embed"
	"fmt"
	"sync"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
)

//go:embed locales/*.yaml
var localesFS embed.FS

// Translator 单例结构体
type Translator struct {
	bundle *i18n.Bundle
	cache  map[string]*i18n.Localizer
}

var (
	instance *Translator
	once     sync.Once
)

// 自定义 YAML 解析函数
func yamlUnmarshal(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}

// GetI18nInstance 获取 Translator 单例
func GetI18nInstance() *Translator {
	once.Do(func() {
		bundle := i18n.NewBundle(language.English)
		bundle.RegisterUnmarshalFunc("yaml", yamlUnmarshal) // 使用 YAML 解析

		instance = &Translator{
			bundle: bundle,
			cache:  make(map[string]*i18n.Localizer), // 语言缓存
		}
	})

	return instance
}

// LoadLanguage 仅加载一次 YAML 资源
func (t *Translator) LoadLanguage(lang string) {
	if _, exists := t.cache[lang]; exists {
		return // 已加载，跳过
	}

	// 从 embed 读取 YAML 文件
	filePath := fmt.Sprintf("locales/%s.yaml", lang)
	_, _ = t.bundle.LoadMessageFileFS(localesFS, filePath)

	t.cache[lang] = i18n.NewLocalizer(t.bundle, lang)
}

// Translate 传入语言和键值获取翻译
func (t *Translator) Translate(lang, key string) string {
	t.LoadLanguage(lang) // 确保语言已加载
	localize := t.cache[lang]
	return localize.MustLocalize(&i18n.LocalizeConfig{
		MessageID: key,
	})
}
