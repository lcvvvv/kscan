package httpfinger

import (
	"regexp"
	"strings"
)

var NewKeywords []string

type keywordFinger []fingerPrint

type fingerPrint struct {
	Cms   string `json:"cms"`
	Rules []Rule `json:"rules"`
}

func (p fingerPrint) match(header string, title string, body string) bool {
	for _, rule := range p.Rules {
		var s string
		if rule.Type == "body" {
			s = body
		}
		if rule.Type == "header" {
			s = header
		}
		if rule.Type == "title" {
			s = title
		}
		if rule.match(s) == false {
			return false
		}
	}
	return true
}

type Rule struct {
	Type    string `json:"type"`
	Keyword string `json:"keyword"`
}

func (r Rule) match(s string) bool {
	s = strings.ToLower(s)
	keyword := strings.ToLower(r.Keyword)
	if keyword[:1] == "~" {
		return regexp.MustCompile(keyword[1:]).MatchString(s)
	} else {
		return strings.Contains(s, keyword)
	}
}

var KeywordFinger keywordFinger

var (
	serverRegx      = regexp.MustCompile(`\n(?i)(server: .*)`)
	xPoweredByRegx  = regexp.MustCompile(`\n(?i)(X-Powered-By: .*)`)
	xRedirectByRegx = regexp.MustCompile(`\n(?i)(X-Redirect-By: .*)`)
)

func (k keywordFinger) Match(header string, title string, body string) string {
	var cmsSlice []string

	for _, kSub := range k {
		if kSub.match(header, title, body) {
			cmsSlice = append(cmsSlice, kSub.Cms)
		}
	}

	keywordSlice := getHeaderDigest(header)

	for _, value := range keywordSlice {
		if isOldKeyword(value) == false {
			NewKeywords = append(NewKeywords, value)
		}
	}

	return strings.Join(cmsSlice, ",")
}

func getHeaderDigest(header string) []string {
	var digest []string

	if serverRegx.MatchString(header) {
		digest = append(digest, serverRegx.FindStringSubmatch(header)[1])
	}

	if xPoweredByRegx.MatchString(header) {
		digest = append(digest, xPoweredByRegx.FindStringSubmatch(header)[1])
	}

	if xRedirectByRegx.MatchString(header) {
		digest = append(digest, xRedirectByRegx.FindStringSubmatch(header)[1])
	}

	return digest
}

func isOldKeyword(value string) bool {
	for _, keyword := range keywordServerSlice {
		value = strings.ToLower(value)
		keyword = "server: " + strings.ToLower(keyword)
		if len(value) < len(keyword) {
			continue
		}
		if value[:len(keyword)] == keyword {
			return true
		}
	}
	for _, keyword := range keywordXPoweredBySlice {
		value = strings.ToLower(value)
		keyword = "x-powered-by: " + strings.ToLower(keyword)
		if len(value) < len(keyword) {
			continue
		}
		if value[:len(keyword)] == keyword {
			return true
		}
	}
	return false
}

var keywordFingerSourceByte = []byte(`
[
	{
		"cms": "ASP.NET",
		"rules": [{
			"type": "header",
			"keyword": "X-Powered-By: ASP.NET"
		}]
	},
	{
		"cms": "天融信",
		"rules": [{
			"type": "header",
			"keyword": "X-Powered-By: topsec"
		}]
	},
	{
		"cms": "PHP",
		"rules": [{
			"type": "header",
			"keyword": "X-Powered-By: PHP"
		}]
	},
	{
		"cms": "ThinkPHP",
		"rules": [{
			"type": "header",
			"keyword": "X-Powered-By: ThinkPHP"
		}]
	},
	{
		"cms": "seeyon",
		"rules": [{
			"type": "body",
			"keyword": "/seeyon/USER-DATA/IMAGES/LOGIN/login.gif"
		}]
	},
	{
		"cms": "seeyon",
		"rules": [{
			"type": "body",
			"keyword": "/seeyon/common/"
		}]
	},
	{
		"cms": "Spring env",
		"rules": [{
			"type": "body",
			"keyword": "servletContextInitParams"
		}]
	},
	{
		"cms": "Spring env",
		"rules": [{
			"type": "body",
			"keyword": "logback"
		}]
	},
	{
		"cms": "Weblogic",
		"rules": [{
			"type": "body",
			"keyword": "Error 404--Not Found"
		}]
	},
	{
		"cms": "Weblogic",
		"rules": [{
			"type": "body",
			"keyword": "Error 403--"
		}]
	},
	{
		"cms": "Weblogic",
		"rules": [{
			"type": "body",
			"keyword": "/console/framework/skins/wlsconsole/images/login_WebLogic_branding.png"
		}]
	},
	{
		"cms": "Weblogic",
		"rules": [{
			"type": "body",
			"keyword": "Welcome to Weblogic Application Server"
		}]
	},
	{
		"cms": "Weblogic",
		"rules": [{
			"type": "body",
			"keyword": "<i>Hypertext Transfer Protocol -- HTTP/1.1</i>"
		}]
	},
	{
		"cms": "Sangfor SSL VPN",
		"rules": [{
			"type": "body",
			"keyword": "/por/login_psw.csp"
		}]
	},
	{
		"cms": "Sangfor SSL VPN",
		"rules": [{
			"type": "body",
			"keyword": "loginPageSP/loginPrivacy.js"
		}]
	},
	{
		"cms": "e-mobile",
		"rules": [{
			"type": "body",
			"keyword": "weaver,e-mobile"
		}]
	},
	{
		"cms": "ecology",
		"rules": [{
			"type": "header",
			"keyword": "ecology_JSessionid"
		}]
	},
	{
		"cms": "Shiro",
		"rules": [{
			"type": "header",
			"keyword": "rememberMe="
		}]
	},
	{
		"cms": "Shiro",
		"rules": [{
			"type": "header",
			"keyword": "=deleteMe"
		}]
	},
	{
		"cms": "e-Bridge",
		"rules": [{
			"type": "body",
			"keyword": "wx.weaver"
		}]
	},
	{
		"cms": "e-Bridge",
		"rules": [{
			"type": "body",
			"keyword": "e-Bridge"
		}]
	},
	{
		"cms": "Swagger UI",
		"rules": [{
			"type": "body",
			"keyword": "Swagger UI"
		}]
	},
	{
		"cms": "Ruijie",
		"rules": [{
			"type": "body",
			"keyword": "4008 111 000"
		}]
	},
	{
		"cms": "Huawei SMC",
		"rules": [{
			"type": "body",
			"keyword": "Script/SmcScript.js?version="
		}]
	},
	{
		"cms": "H3C Router",
		"rules": [{
			"type": "body",
			"keyword": "/wnm/ssl/web/frame/login.html"
		}]
	},
	{
		"cms": "Cisco SSLVPN",
		"rules": [{
			"type": "body",
			"keyword": "/+CSCOE+/logon.html"
		}]
	},
	{
		"cms": "\u901a\u8fbeOA",
		"rules": [{
			"type": "body",
			"keyword": "/images/tongda.ico"
		}]
	},
	{
		"cms": "\u901a\u8fbeOA",
		"rules": [{
			"type": "body",
			"keyword": "Office Anywhere"
		}]
	},
	{
		"cms": "\u901a\u8fbeOA",
		"rules": [{
			"type": "body",
			"keyword": "\u901a\u8fbeOA"
		}]
	},
	{
		"cms": "\u6df1\u4fe1\u670d waf",
		"rules": [{
			"type": "body",
			"keyword": "rsa.js"
		}]
	},
	{
		"cms": "\u6df1\u4fe1\u670d waf",
		"rules": [{
			"type": "body",
			"keyword": "Redirect to..."
		}]
	},
	{
		"cms": "\u7f51\u5fa1 vpn",
		"rules": [{
			"type": "body",
			"keyword": "/vpn/common/js/leadsec.js"
		}]
	},
	{
		"cms": "\u542f\u660e\u661f\u8fb0\u5929\u6e05\u6c49\u9a6cUSG\u9632\u706b\u5899",
		"rules": [{
			"type": "body",
			"keyword": "/cgi-bin/webui?op=get_product_model"
		}]
	},
	{
		"cms": "\u84dd\u51cc OA",
		"rules": [{
			"type": "body",
			"keyword": "sys/ui/extend/theme/default/style/icon.css"
		}]
	},
	{
		"cms": "\u6df1\u4fe1\u670d\u4e0a\u7f51\u884c\u4e3a\u7ba1\u7406\u7cfb\u7edf",
		"rules": [{
			"type": "body",
			"keyword": "utccjfaewjb = function(str, key)"
		}]
	},
	{
		"cms": "\u6df1\u4fe1\u670d\u4e0a\u7f51\u884c\u4e3a\u7ba1\u7406\u7cfb\u7edf",
		"rules": [{
			"type": "body",
			"keyword": "document.write(WRFWWCSFBXMIGKRKHXFJ"
		}]
	},
	{
		"cms": "\u6df1\u4fe1\u670d\u5e94\u7528\u4ea4\u4ed8\u62a5\u8868\u7cfb\u7edf",
		"rules": [{
			"type": "body",
			"keyword": "/reportCenter/index.php?cls_mode=cluster_mode_others"
		}]
	},
	{
		"cms": "\u91d1\u8776\u4e91\u661f\u7a7a",
		"rules": [{
			"type": "body",
			"keyword": "HTML5/content/themes/kdcss.min.css"
		}]
	},
	{
		"cms": "\u91d1\u8776\u4e91\u661f\u7a7a",
		"rules": [{
			"type": "body",
			"keyword": "/ClientBin/Kingdee.BOS.XPF.App.xap"
		}]
	},
	{
		"cms": "CoreMail",
		"rules": [{
			"type": "body",
			"keyword": "coremail/common"
		}]
	},
	{
		"cms": "\u542f\u660e\u661f\u8fb0\u5929\u6e05\u6c49\u9a6cUSG\u9632\u706b\u5899",
		"rules": [{
			"type": "body",
			"keyword": "\u5929\u6e05\u6c49\u9a6cUSG"
		}]
	},
	{
		"cms": "Jboss",
		"rules": [{
			"type": "body",
			"keyword": "jboss.css"
		}]
	},
	{
		"cms": "Gitlab",
		"rules": [{
			"type": "body",
			"keyword": "assets/gitlab_logo"
		}]
	},
	{
		"cms": "\u5b9d\u5854-BT.cn",
		"rules": [{
			"type": "body",
			"keyword": "\u5165\u53e3\u6821\u9a8c\u5931\u8d25"
		}]
	},
	{
		"cms": "\u7985\u9053",
		"rules": [{
			"type": "body",
			"keyword": "/theme/default/images/main/zt-logo.png"
		}]
	},
	{
		"cms": "ADSL/Router",
		"rules": [{
			"type": "body",
			"keyword": "zentaosid"
		}]
	},
	{
		"cms": "ADSL/Router",
		"rules": [{
			"type": "header",
			"keyword": "zentaosid"
		}]
	},
	{
		"cms": "\u7528\u53cb\u8f6f\u4ef6",
		"rules": [{
			"type": "body",
			"keyword": "UFIDA Software CO.LTD all rights reserved"
		}]
	},
	{
		"cms": "YONYOU NC",
		"rules": [{
			"type": "body",
			"keyword": "uclient.yonyou.com"
		}]
	},
	{
		"cms": "\u5b9d\u5854-BT.cn",
		"rules": [{
			"type": "body",
			"keyword": "\u5b9d\u5854Linux\u9762\u677f"
		}]
	},
	{
		"cms": "RabbitMQ",
		"rules": [{
			"type": "body",
			"keyword": "<title>RabbitMQ Management</title>"
		}]
	},
	{
		"cms": "Zabbix",
		"rules": [{
			"type": "body",
			"keyword": "zabbix"
		}]
	},
	{
		"cms": "\u8054\u8f6f\u51c6\u5165",
		"rules": [{
			"type": "body",
			"keyword": "\u7f51\u7edc\u51c6\u5165"
		}]
	},
	{
		"cms": "\u5217\u76ee\u5f55",
		"rules": [{
			"type": "body",
			"keyword": "Index of /"
		}]
	},
	{
		"cms": "\u5217\u76ee\u5f55",
		"rules": [{
			"type": "body",
			"keyword": " - /</title>"
		}]
	},
	{
		"cms": "\u6d6a\u6f6e\u670d\u52a1\u5668IPMI\u7ba1\u7406\u53e3",
		"rules": [{
			"type": "body",
			"keyword": "img/inspur_logo.png"
		}]
	},
	{
		"cms": "RegentApi_v2.0",
		"rules": [{
			"type": "body",
			"keyword": "RegentApi_v2.0"
		}]
	},
	{
		"cms": "Tomcat\u9ed8\u8ba4\u9875\u9762",
		"rules": [{
			"type": "body",
			"keyword": "/manager/status"
		}]
	},
	{
	"cms": "Verizon_Wireless_Router",
	"rules": [{
		"type": "title",
		"keyword": "Wireless Broadband Router Management Console"
	}, {
		"type": "body",
		"keyword": "verizon_logo_blk.gif"
	}]
}, {
	"cms": "NSFOCUS_WAF",
	"rules": [{
		"type": "title",
		"keyword": "WAF NSFOCUS"
	}, {
		"type": "body",
		"keyword": "/images/logo/nsfocus.png"
	}]
}, {
	"cms": "IndusGuard_WAF",
	"rules": [{
		"type": "title",
		"keyword": "IndusGuard WAF"
	}, {
		"type": "body",
		"keyword": "wafportal/wafportal.nocache.js"
	}]
}, {
	"cms": "Maticsoft_Shop_动软商城",
	"rules": [{
		"type": "body",
		"keyword": "maticsoft"
	}, {
		"type": "body",
		"keyword": "/Areas/Shop/"
	}]
}, {
	"cms": "MaticsoftSNS_动软分享社区",
	"rules": [{
		"type": "body",
		"keyword": "maticsoft"
	}, {
		"type": "body",
		"keyword": "/Areas/SNS/"
	}]
}, {
	"cms": "梭子鱼防火墙",
	"rules": [{
		"type": "body",
		"keyword": "http://www.barracudanetworks.com?a=bsf_product\" class=\"transbutton"
	}, {
		"type": "body",
		"keyword": "/cgi-mod/header_logo.cgi"
	}]
}, {
	"cms": "twcms",
	"rules": [{
		"type": "body",
		"keyword": "/twcms/theme/"
	}, {
		"type": "body",
		"keyword": "/css/global.css"
	}]
}, {
	"cms": "QNO_Router",
	"rules": [{
		"type": "body",
		"keyword": "/QNOVirtual_Keyboard.js"
	}, {
		"type": "body",
		"keyword": "/images/login_img01_03.gif"
	}]
}, {
	"cms": "TerraMaster",
	"rules": [{
		"type": "title",
		"keyword": "TerraMaster"
	}, {
		"type": "body",
		"keyword": "/js/common.js"
	}]
}, {
	"cms": "sdcms",
	"rules": [{
		"type": "body",
		"keyword": "var webroot="
	}, {
		"type": "body",
		"keyword": "/js/sdcms.js"
	}]
}, {
	"cms": "Joomla",
	"rules": [{
		"type": "body",
		"keyword": "/media/system/js/core.js"
	}, {
		"type": "body",
		"keyword": "/media/system/js/mootools-core.js"
	}]
}, {
	"cms": "爱快流控路由",
	"rules": [{
		"type": "title",
		"keyword": "爱快"
	}, {
		"type": "body",
		"keyword": "/resources/images/land_prompt_ico01.gif"
	}]
}, {
	"cms": "ESPCMS",
	"rules": [{
		"type": "body",
		"keyword": "infolist_fff"
	}, {
		"type": "body",
		"keyword": "/templates/default/style/tempates_div.css"
	}]
}, {
	"cms": "百为智能流控路由器",
	"rules": [{
		"type": "title",
		"keyword": "BYTEVALUE 智能流控路由器"
	}, {
		"type": "body",
		"keyword": "<a href=\"http://www.bytevalue.com/\" target=\"_blank\">"
	}]
}, {
	"cms": "UBNT_UniFi系列路由",
	"rules": [{
		"type": "title",
		"keyword": "UniFi"
	}, {
		"type": "body",
		"keyword": "<div class=\"appGlobalHeader\">"
	}]
}, {
	"cms": "乐视路由器",
	"rules": [{
		"type": "title",
		"keyword": "乐视路由器"
	}, {
		"type": "body",
		"keyword": "<div class=\"login-logo\"></div>"
	}]
}, {
	"cms": "斐讯Fortress",
	"rules": [{
		"type": "title",
		"keyword": "斐讯Fortress防火墙"
	}, {
		"type": "body",
		"keyword": "<meta name=\"author\" content=\"上海斐讯数据通信技术有限公司\" />"
	}]
}, {
	"cms": "68 Classifieds",
	"rules": [{
		"type": "body",
		"keyword": "powered by"
	}, {
		"type": "body",
		"keyword": "68 Classifieds"
	}]
}, {
	"cms": "Aardvark Topsites",
	"rules": [{
		"type": "body",
		"keyword": "Powered by"
	}, {
		"type": "body",
		"keyword": "Aardvark Topsites"
	}]
}, {
	"cms": "Adiscon_LogAnalyzer",
	"rules": [{
		"type": "body",
		"keyword": "Adiscon LogAnalyzer"
	}, {
		"type": "body",
		"keyword": "Adiscon GmbH"
	}]
}, {
	"cms": "AllNewsManager_NET",
	"rules": [{
		"type": "body",
		"keyword": "Powered by"
	}, {
		"type": "body",
		"keyword": "AllNewsManager"
	}]
}, {
	"cms": "ARRIS-Touchstone-Router",
	"rules": [{
		"type": "body",
		"keyword": "Copyright"
	}, {
		"type": "body",
		"keyword": "ARRIS Group"
	}]
}, {
	"cms": "Aruba-Device",
	"rules": [{
		"type": "body",
		"keyword": "Copyright"
	}, {
		"type": "body",
		"keyword": "Aruba Networks"
	}]
}, {
	"cms": "ashnews",
	"rules": [{
		"type": "body",
		"keyword": "powered by"
	}, {
		"type": "body",
		"keyword": "ashnews"
	}]
}, {
	"cms": "Atomic-Photo-Album",
	"rules": [{
		"type": "body",
		"keyword": "Powered by"
	}, {
		"type": "body",
		"keyword": "Atomic Photo Album"
	}]
}, {
	"cms": "Redmine",
	"rules": [{
		"type": "body",
		"keyword": "Redmine"
	}, {
		"type": "body",
		"keyword": "authenticity_token"
	}]
}, {
	"cms": "beecms",
	"rules": [{
		"type": "body",
		"keyword": "powerd by"
	}, {
		"type": "body",
		"keyword": "BEESCMS"
	}]
}, {
	"cms": "Magento",
	"rules": [{
		"type": "body",
		"keyword": "/skin/frontend/"
	}, {
		"type": "body",
		"keyword": "BLANK_IMG"
	}]
}, {
	"cms": "Emlog-PHP",
	"rules": [{
		"type": "body",
		"keyword": "/include/lib/js/common_tpl.js"
	}, {
		"type": "body",
		"keyword": "content/templates"
	}]
}, {
	"cms": "OA企业智能办公自动化系统",
	"rules": [{
		"type": "body",
		"keyword": "input name=\"S1\" type=\"image\""
	}, {
		"type": "body",
		"keyword": "count/mystat.asp"
	}]
}, {
	"cms": "CISCO_EPC3925",
	"rules": [{
		"type": "body",
		"keyword": "Docsis_system"
	}, {
		"type": "body",
		"keyword": "EPC3925"
	}]
}, {
	"cms": "地平线CMS",
	"rules": [{
		"type": "body",
		"keyword": "search_result.aspx"
	}, {
		"type": "body",
		"keyword": "frmsearch"
	}]
}, {
	"cms": "国家数字化学习资源中心系统",
	"rules": [{
		"type": "title",
		"keyword": "页面加载中,请稍候"
	}, {
		"type": "body",
		"keyword": "FrontEnd"
	}]
}, {
	"cms": "H3C公司产品",
	"rules": [{
		"type": "body",
		"keyword": "Copyright"
	}, {
		"type": "body",
		"keyword": "H3C Corporation"
	}]
}, {
	"cms": "BlognPlus",
	"rules": [{
		"type": "body",
		"keyword": "Powered by"
	}, {
		"type": "body",
		"keyword": "href=\"http://www.blogn.org"
	}]
}, {
	"cms": "Schneider_Quantum_140NOE77101",
	"rules": [{
		"type": "body",
		"keyword": "indexLanguage"
	}, {
		"type": "body",
		"keyword": "html/config.js"
	}]
}, {
	"cms": "AlstraSoft-AskMe",
	"rules": [{
		"type": "body",
		"keyword": "Powered by"
	}, {
		"type": "body",
		"keyword": "http://www.alstrasoft.com"
	}]
}, {
	"cms": "BlogEngine_NET",
	"rules": [{
		"type": "body",
		"keyword": "Powered by"
	}, {
		"type": "body",
		"keyword": "http://www.dotnetblogengine.net"
	}]
}, {
	"cms": "iGaming-CMS",
	"rules": [{
		"type": "body",
		"keyword": "Powered by"
	}, {
		"type": "body",
		"keyword": "http://www.igamingcms.com/"
	}]
}, {
	"cms": "jcg无线路由器",
	"rules": [{
		"type": "title",
		"keyword": "Wireless Router"
	}, {
		"type": "body",
		"keyword": "http://www.jcgcn.com"
	}]
}, {
	"cms": "jcg无线路由器",
	"rules": [{
		"type": "title",
		"keyword": "Wireless Router"
	}, {
		"type": "body",
		"keyword": "http://www.jcgcn.com"
	}]
}, {
	"cms": "jobberBase",
	"rules": [{
		"type": "body",
		"keyword": "powered by"
	}, {
		"type": "body",
		"keyword": "http://www.jobberbase.com"
	}]
}, {
	"cms": "PhpCMS",
	"rules": [{
		"type": "body",
		"keyword": "Powered by"
	}, {
		"type": "body",
		"keyword": "http://www.phpcms.cn"
	}]
}, {
	"cms": "TCCMS",
	"rules": [{
		"type": "body",
		"keyword": "index.php?ac=link_more"
	}, {
		"type": "body",
		"keyword": "index.php?ac=news_list"
	}]
}, {
	"cms": "逐浪zoomla",
	"rules": [{
		"type": "body",
		"keyword": "NodePage.aspx"
	}, {
		"type": "body",
		"keyword": "Item"
	}]
}, {
	"cms": "海康威视iVMS",
	"rules": [{
		"type": "body",
		"keyword": "g_szCacheTime"
	}, {
		"type": "body",
		"keyword": "iVMS"
	}]
}, {
	"cms": "rcms",
	"rules": [{
		"type": "body",
		"keyword": "/r/cms/www/"
	}, {
		"type": "body",
		"keyword": "jhtml"
	}]
}, {
	"cms": "Polycom",
	"rules": [{
		"type": "title",
		"keyword": "Polycom"
	}, {
		"type": "body",
		"keyword": "kAllowDirectHTMLFileAccess"
	}]
}, {
	"cms": "用友商战实践平台",
	"rules": [{
		"type": "body",
		"keyword": "Login_Main_BG"
	}, {
		"type": "body",
		"keyword": "Login_Owner"
	}]
}, {
	"cms": "EasyTrace(botwave)",
	"rules": [{
		"type": "title",
		"keyword": "EasyTrace"
	}, {
		"type": "body",
		"keyword": "login_page"
	}]
}, {
	"cms": "北京清科锐华CEMIS",
	"rules": [{
		"type": "body",
		"keyword": "/theme/2009/image"
	}, {
		"type": "body",
		"keyword": "login.asp"
	}]
}, {
	"cms": "北京清科锐华CEMIS",
	"rules": [{
		"type": "body",
		"keyword": "/theme/2009/image"
	}, {
		"type": "body",
		"keyword": "login.asp"
	}]
}, {
	"cms": "UFIDA_NC",
	"rules": [{
		"type": "body",
		"keyword": "UFIDA"
	}, {
		"type": "body",
		"keyword": "logo/images/"
	}]
}, {
	"cms": "UFIDA_NC",
	"rules": [{
		"type": "body",
		"keyword": "UFIDA"
	}, {
		"type": "body",
		"keyword": "logo/images/"
	}]
}, {
	"cms": "UFIDA_NC",
	"rules": [{
		"type": "body",
		"keyword": "UFIDA"
	}, {
		"type": "body",
		"keyword": "logo/images/"
	}]
}, {
	"cms": "HIMS酒店云计算服务",
	"rules": [{
		"type": "body",
		"keyword": "GB_ROOT_DIR"
	}, {
		"type": "body",
		"keyword": "maincontent.css"
	}]
}, {
	"cms": "mikrotik",
	"rules": [{
		"type": "title",
		"keyword": "RouterOS"
	}, {
		"type": "body",
		"keyword": "mikrotik"
	}]
}, {
	"cms": "mikrotik",
	"rules": [{
		"type": "title",
		"keyword": "RouterOS"
	}, {
		"type": "body",
		"keyword": "mikrotik"
	}]
}, {
	"cms": "管理易",
	"rules": [{
		"type": "body",
		"keyword": "管理易"
	}, {
		"type": "body",
		"keyword": "minierp"
	}]
}, {
	"cms": "h3c路由器",
	"rules": [{
		"type": "title",
		"keyword": "Web user login"
	}, {
		"type": "body",
		"keyword": "nLanguageSupported"
	}]
}, {
	"cms": "h3c路由器",
	"rules": [{
		"type": "title",
		"keyword": "Web user login"
	}, {
		"type": "body",
		"keyword": "nLanguageSupported"
	}]
}, {
	"cms": "OpenSNS",
	"rules": [{
		"type": "body",
		"keyword": "powered by"
	}, {
		"type": "body",
		"keyword": "opensns"
	}]
}, {
	"cms": "惠尔顿上网行为管理系统",
	"rules": [{
		"type": "body",
		"keyword": "updateLoginPswd.php"
	}, {
		"type": "body",
		"keyword": "PassroedEle"
	}]
}, {
	"cms": "discuz",
	"rules": [{
		"type": "body",
		"keyword": "discuz_uid"
	}, {
		"type": "body",
		"keyword": "portal.php?mod=view"
	}]
}, {
	"cms": "74cms",
	"rules": [{
		"type": "body",
		"keyword": "/templates/default/css/common.css"
	}, {
		"type": "body",
		"keyword": "selectjobscategory"
	}]
}, {
	"cms": "SiteServer",
	"rules": [{
		"type": "body",
		"keyword": "siteserver"
	}, {
		"type": "body",
		"keyword": "sitefiles"
	}]
}, {
	"cms": "TeamViewer",
	"rules": [{
		"type": "body",
		"keyword": "This site is running"
	}, {
		"type": "body",
		"keyword": "TeamViewer"
	}]
}, {
	"cms": "Typecho",
	"rules": [{
		"type": "body",
		"keyword": "强力驱动"
	}, {
		"type": "body",
		"keyword": "Typecho"
	}]
}, {
	"cms": "yongyoufe",
	"rules": [{
		"type": "body",
		"keyword": "V_show"
	}, {
		"type": "body",
		"keyword": "V_hedden"
	}]
}, {
	"cms": "创星伟业校园网群",
	"rules": [{
		"type": "body",
		"keyword": "javascripts/float.js"
	}, {
		"type": "body",
		"keyword": "vcxvcxv"
	}]
}, {
	"cms": "创星伟业校园网群",
	"rules": [{
		"type": "body",
		"keyword": "javascripts/float.js"
	}, {
		"type": "body",
		"keyword": "vcxvcxv"
	}]
}, {
	"cms": "phpinfo",
	"rules": [{
		"type": "title",
		"keyword": "phpinfo"
	}, {
		"type": "body",
		"keyword": "Virtual Directory Support "
	}]
}, {
	"cms": "trs_wcm",
	"rules": [{
		"type": "body",
		"keyword": "forum.trs.com.cn"
	}, {
		"type": "body",
		"keyword": "wcm"
	}]
}, {
	"cms": "Advanced-Image-Hosting-Script",
	"rules": [{
		"type": "body",
		"keyword": "Powered by"
	}, {
		"type": "body",
		"keyword": "yabsoft.com\" "
	}]
}, {
	"cms": "360webfacil_360WebManager",
	"rules": [{
		"type": "body",
		"keyword": "publico/template/"
	}, {
		"type": "body",
		"keyword": "zonapie"
	}]
}, {
	"cms": "ZTE_ZSRV2_Router",
	"rules": [{
		"type": "title",
		"keyword": "ZSRV2路由器Web管理系统"
	}, {
		"type": "body",
		"keyword": "ZTE Corporation. All Rights Reserved."
	}]
}, {
	"cms": "webplus",
	"rules": [{
		"type": "body",
		"keyword": "webplus"
	}, {
		"type": "body",
		"keyword": "高校网站群管理平台"
	}]
}, {
	"cms": "育友软件",
	"rules": [{
		"type": "body",
		"keyword": "http://www.yuysoft.com/"
	}, {
		"type": "body",
		"keyword": "技术支持"
	}]
}, {
	"cms": "通达OA",
	"rules": [{
		"type": "body",
		"keyword": "OA提示：不能登录OA"
	}, {
		"type": "body",
		"keyword": "紧急通知：今日10点停电"
	}]
}, {
	"cms": "ThinkSNS",
	"rules": [{
		"type": "body",
		"keyword": "/addons/theme/"
	}, {
		"type": "body",
		"keyword": "全局变量"
	}]
}, {
	"cms": "合正网站群内容管理系统",
	"rules": [{
		"type": "body",
		"keyword": "Produced By"
	}, {
		"type": "body",
		"keyword": "网站群内容管理系统"
	}]
}, {
	"cms": "pmway_E4_crm",
	"rules": [{
		"type": "title",
		"keyword": "E4"
	}, {
		"type": "title",
		"keyword": "CRM"
	}]
}, {
	"cms": "pmway_E4_crm",
	"rules": [{
		"type": "title",
		"keyword": "E4"
	}, {
		"type": "title",
		"keyword": "CRM"
	}]
}, {
	"cms": "Everything",
	"rules": [{
		"type": "body",
		"keyword": "everything.png\")"
	}, {
		"type": "title",
		"keyword": "Everything"
	}]
}, {
	"cms": "Nexus_NX_router",
	"rules": [{
		"type": "body",
		"keyword": "http://nexuswifi.com/"
	}, {
		"type": "title",
		"keyword": "Nexus NX"
	}]
}, {
	"cms": "NetShare_VPN",
	"rules": [{
		"type": "title",
		"keyword": "NetShare"
	}, {
		"type": "title",
		"keyword": "VPN"
	}]
}, {
	"cms": "网御WAF",
	"rules": [{
		"type": "body",
		"keyword": "<div id=\"divLogin\">"
	}, {
		"type": "title",
		"keyword": "网御WAF"
	}]
}, {
	"cms": "网易企业邮箱",
	"rules": [{
		"type": "body",
		"keyword": "frmvalidator"
	}, {
		"type": "title",
		"keyword": "邮箱用户登录"
	}]
}, {
	"cms": "CEMIS",
	"rules": [{
		"type": "body",
		"keyword": "<div id=\"demo\" style=\"overflow:hidden"
	}, {
		"type": "title",
		"keyword": "综合项目管理系统登录"
	}]
}
]
`)

var keywordFingerFofaByte = []byte(`
[{
		"cms": "Everything",
		"rules": [{
			"type": "body",
			"keyword": "Everything.gif"
		}]
	},
	{
		"cms": "通达OA",
		"rules": [{
			"type": "body",
			"keyword": "<a href='http://www.tongda2000.com/' target='_black'>通达官网</a></div>"
		}]
	},
	{
		"cms": "TurboCMS",
		"rules": [{
			"type": "body",
			"keyword": "/cmsapp/zxdcADD.jsp"
		}]
	},
	{
		"cms": "协众OA",
		"rules": [{
			"type": "body",
			"keyword": "Powered by 协众OA"
		}]
	},
	{
		"cms": "Django",
		"rules": [{
			"type": "body",
			"keyword": "__admin_media_prefix__"
		}]
	},
	{
		"cms": "Adobe_ CQ5",
		"rules": [{
			"type": "body",
			"keyword": "_jcr_content"
		}]
	},
	{
		"cms": "ZCMS",
		"rules": [{
			"type": "body",
			"keyword": "_ZCMS_ShowNewMessage"
		}]
	},
	{
		"cms": "xheditor",
		"rules": [{
			"type": "body",
			"keyword": ".xheditor("
		}]
	},
	{
		"cms": "shopex",
		"rules": [{
			"type": "body",
			"keyword": "@author litie[aita]shopex.cn"
		}]
	},
	{
		"cms": "eMeeting-Online-Dating-Software",
		"rules": [{
			"type": "body",
			"keyword": "/_eMeetingGlobals.js"
		}]
	},
	{
		"cms": "迈捷邮件系统(MagicMail)",
		"rules": [{
			"type": "body",
			"keyword": "/aboutus/magicmail.gif"
		}]
	},
	{
		"cms": "Astaro-Command-Center",
		"rules": [{
			"type": "body",
			"keyword": "/acc_aggregated_reporting.js"
		}]
	},
	{
		"cms": "PineApp",
		"rules": [{
			"type": "body",
			"keyword": "/admin/css/images/pineapp.ico"
		}]
	},
	{
		"cms": "AnyGate",
		"rules": [{
			"type": "body",
			"keyword": "/anygate.php"
		}]
	},
	{
		"cms": "中望OA",
		"rules": [{
			"type": "body",
			"keyword": "/app_qjuserinfo/qjuserinfoadd.jsp"
		}]
	},
	{
		"cms": "ThinkSAAS",
		"rules": [{
			"type": "body",
			"keyword": "/app/home/skins/default/style.css"
		}]
	},
	{
		"cms": "Apache-Archiva",
		"rules": [{
			"type": "body",
			"keyword": "/archiva.css"
		}]
	},
	{
		"cms": "Apache-Archiva",
		"rules": [{
			"type": "body",
			"keyword": "/archiva.js"
		}]
	},
	{
		"cms": "ARRIS-Touchstone-Router",
		"rules": [{
			"type": "body",
			"keyword": "/arris_style.css"
		}]
	},
	{
		"cms": "Aurion",
		"rules": [{
			"type": "body",
			"keyword": "/aurion.js"
		}]
	},
	{
		"cms": "易瑞授权访问系统",
		"rules": [{
			"type": "body",
			"keyword": "/authjsp/login.jsp"
		}]
	},
	{
		"cms": "Barracuda-Spam-Firewall",
		"rules": [{
			"type": "body",
			"keyword": "/barracuda.css"
		}]
	},
	{
		"cms": "Biscom-Delivery-Server",
		"rules": [{
			"type": "body",
			"keyword": "/bds/includes/fdsJavascript.do"
		}]
	},
	{
		"cms": "Biscom-Delivery-Server",
		"rules": [{
			"type": "body",
			"keyword": "/bds/stylesheets/fds.css"
		}]
	},
	{
		"cms": "Advantech-WebAccess",
		"rules": [{
			"type": "body",
			"keyword": "/broadWeb/bwuconfig.asp"
		}]
	},
	{
		"cms": "Advantech-WebAccess",
		"rules": [{
			"type": "body",
			"keyword": "/broadweb/WebAccessClientSetup.exe"
		}]
	},
	{
		"cms": "orocrm",
		"rules": [{
			"type": "body",
			"keyword": "/bundles/oroui/"
		}]
	},
	{
		"cms": "orocrm",
		"rules": [{
			"type": "body",
			"keyword": "/bundles/oroui/"
		}]
	},
	{
		"cms": "Advantech-WebAccess",
		"rules": [{
			"type": "body",
			"keyword": "/bw_templete1.dwt"
		}]
	},
	{
		"cms": "CafeEngine",
		"rules": [{
			"type": "body",
			"keyword": "/CafeEngine/style.css"
		}]
	},
	{
		"cms": "cApexWEB",
		"rules": [{
			"type": "body",
			"keyword": "/capexweb.parentvalidatepassword"
		}]
	},
	{
		"cms": "cisco UCM",
		"rules": [{
			"type": "body",
			"keyword": "/ccmadmin/"
		}]
	},
	{
		"cms": "BlueNet-Video",
		"rules": [{
			"type": "body",
			"keyword": "/cgi-bin/client_execute.cgi?tUD=0"
		}]
	},
	{
		"cms": "IBM-Cognos",
		"rules": [{
			"type": "body",
			"keyword": "/cgi-bin/cognos.cgi"
		}]
	},
	{
		"cms": "Axis-PrintServer",
		"rules": [{
			"type": "body",
			"keyword": "/cgi-bin/prodhelp?prod="
		}]
	},
	{
		"cms": "GenOHM-SCADA",
		"rules": [{
			"type": "body",
			"keyword": "/cgi-bin/scada-vis/"
		}]
	},
	{
		"cms": "Spammark邮件信息安全网关",
		"rules": [{
			"type": "body",
			"keyword": "/cgi-bin/spammark?empty=1"
		}]
	},
	{
		"cms": "TurboCMS",
		"rules": [{
			"type": "body",
			"keyword": "/cmsapp/count/newstop_index.jsp?siteid="
		}]
	},
	{
		"cms": "sugon_gridview",
		"rules": [{
			"type": "body",
			"keyword": "/common/resources/images/common/app/gridview.ico"
		}]
	},
	{
		"cms": "Acidcat CMS",
		"rules": [{
			"type": "body",
			"keyword": "/css/admin_import.css"
		}]
	},
	{
		"cms": "CMSTop",
		"rules": [{
			"type": "body",
			"keyword": "/css/cmstop-common.css"
		}]
	},
	{
		"cms": "锐捷 RG-DBS",
		"rules": [{
			"type": "body",
			"keyword": "/css/impl-security.css"
		}]
	},
	{
		"cms": "mymps",
		"rules": [{
			"type": "body",
			"keyword": "/css/mymps.css"
		}]
	},
	{
		"cms": "DMXReady-Portfolio-Manager",
		"rules": [{
			"type": "body",
			"keyword": "/css/PortfolioManager/styles_display_page.css"
		}]
	},
	{
		"cms": "weiphp",
		"rules": [{
			"type": "body",
			"keyword": "/css/weiphp.css"
		}]
	},
	{
		"cms": "Yxcms",
		"rules": [{
			"type": "body",
			"keyword": "/css/yxcms.css"
		}]
	},
	{
		"cms": "一采通",
		"rules": [{
			"type": "body",
			"keyword": "/custom/GroupNewsList.aspx?GroupId="
		}]
	},
	{
		"cms": "一采通",
		"rules": [{
			"type": "body",
			"keyword": "/custom/GroupNewsList.aspx?GroupId="
		}]
	},
	{
		"cms": "Juniper-NetScreen-Secure-Access",
		"rules": [{
			"type": "body",
			"keyword": "/dana-na/auth/welcome.cgi"
		}]
	},
	{
		"cms": "锐捷 RG-DBS",
		"rules": [{
			"type": "body",
			"keyword": "/dbaudit/authenticate"
		}]
	},
	{
		"cms": "ezOFFICE",
		"rules": [{
			"type": "body",
			"keyword": "/defaultroot/js/cookie.js"
		}]
	},
	{
		"cms": "某通用型政府cms",
		"rules": [{
			"type": "body",
			"keyword": "/deptWebsiteAction.do"
		}]
	},
	{
		"cms": "Astaro-Security-Gateway",
		"rules": [{
			"type": "body",
			"keyword": "/doc/astaro-license.txt"
		}]
	},
	{
		"cms": "Donations-Cloud",
		"rules": [{
			"type": "body",
			"keyword": "/donationscloud.css"
		}]
	},
	{
		"cms": "DotCMS",
		"rules": [{
			"type": "body",
			"keyword": "/dotAsset/"
		}]
	},
	{
		"cms": "DrugPak",
		"rules": [{
			"type": "body",
			"keyword": "/dplimg/DPSTYLE.CSS"
		}]
	},
	{
		"cms": "dwr",
		"rules": [{
			"type": "body",
			"keyword": "/dwr/engine.js"
		}]
	},
	{
		"cms": "泛普建筑工程施工OA",
		"rules": [{
			"type": "body",
			"keyword": "/dwr/interface/LoginService.js"
		}]
	},
	{
		"cms": "PageAdmin",
		"rules": [{
			"type": "body",
			"keyword": "/e/images/favicon.ico"
		}]
	},
	{
		"cms": "Echo",
		"rules": [{
			"type": "body",
			"keyword": "/Echo2/echoweb/login"
		}]
	},
	{
		"cms": "eTicket",
		"rules": [{
			"type": "body",
			"keyword": "/eticket/eticket.css"
		}]
	},
	{
		"cms": "ewebeditor",
		"rules": [{
			"type": "body",
			"keyword": "/ewebeditor.htm?"
		}]
	},
	{
		"cms": "fangmail",
		"rules": [{
			"type": "body",
			"keyword": "/fangmail/default/css/em_css.css"
		}]
	},
	{
		"cms": "FormMail",
		"rules": [{
			"type": "body",
			"keyword": "/FormMail.pl"
		}]
	},
	{
		"cms": "Gallery",
		"rules": [{
			"type": "body",
			"keyword": "/gallery/images/gallery.png"
		}]
	},
	{
		"cms": "希尔OA",
		"rules": [{
			"type": "body",
			"keyword": "/heeroa/login.do"
		}]
	},
	{
		"cms": "Hiki",
		"rules": [{
			"type": "body",
			"keyword": "/hiki_base.css"
		}]
	},
	{
		"cms": "Kloxo-Single-Server",
		"rules": [{
			"type": "body",
			"keyword": "/htmllib/js/preop.js"
		}]
	},
	{
		"cms": "元年财务软件",
		"rules": [{
			"type": "body",
			"keyword": "/image/logo/yuannian.gif"
		}]
	},
	{
		"cms": "Aruba-Device",
		"rules": [{
			"type": "body",
			"keyword": "/images/arubalogo.gif"
		}]
	},
	{
		"cms": "Carrier-CCNWeb",
		"rules": [{
			"type": "body",
			"keyword": "/images/CCNWeb.gif"
		}]
	},
	{
		"cms": "Cogent-DataHub",
		"rules": [{
			"type": "body",
			"keyword": "/images/Cogent.gif"
		}]
	},
	{
		"cms": "MetInfo",
		"rules": [{
			"type": "body",
			"keyword": "/images/css/metinfo.css"
		}]
	},
	{
		"cms": "中望OA",
		"rules": [{
			"type": "body",
			"keyword": "/IMAGES/default/first/xtoa_logo.png"
		}]
	},
	{
		"cms": "NITC",
		"rules": [{
			"type": "body",
			"keyword": "/images/nitc1.png"
		}]
	},
	{
		"cms": "SCADA PLC",
		"rules": [{
			"type": "body",
			"keyword": "/images/rockcolor.gif"
		}]
	},
	{
		"cms": "FreeNAS",
		"rules": [{
			"type": "body",
			"keyword": "/images/ui/freenas-logo.png"
		}]
	},
	{
		"cms": "ASPCMS",
		"rules": [{
			"type": "body",
			"keyword": "/inc/AspCms_AdvJs.asp"
		}]
	},
	{
		"cms": "Axis-Network-Camera",
		"rules": [{
			"type": "body",
			"keyword": "/incl/trash.shtml"
		}]
	},
	{
		"cms": "DotCMS",
		"rules": [{
			"type": "body",
			"keyword": "/index.dot"
		}]
	},
	{
		"cms": "PhpCMS",
		"rules": [{
			"type": "body",
			"keyword": "/index.php?m=content&amp;c=index&amp;a=lists"
		}]
	},
	{
		"cms": "PhpCMS",
		"rules": [{
			"type": "body",
			"keyword": "/index.php?m=content&c=index&a=lists"
		}]
	},
	{
		"cms": "O2OCMS",
		"rules": [{
			"type": "body",
			"keyword": "/index.php/clasify/showone/gtitle/"
		}]
	},
	{
		"cms": "Atmail-WebMail",
		"rules": [{
			"type": "body",
			"keyword": "/index.php/mail/auth/processlogin"
		}]
	},
	{
		"cms": "TPshop",
		"rules": [{
			"type": "body",
			"keyword": "/index.php/Mobile/Index/index.html"
		}]
	},
	{
		"cms": "EPiServer",
		"rules": [{
			"type": "body",
			"keyword": "/javascript/episerverscriptmanager.js"
		}]
	},
	{
		"cms": "金笛邮件系统",
		"rules": [{
			"type": "body",
			"keyword": "/jdwm/cgi/login.cgi?login"
		}]
	},
	{
		"cms": "青果软件",
		"rules": [{
			"type": "body",
			"keyword": "/jkingo.js"
		}]
	},
	{
		"cms": "iWebSNS",
		"rules": [{
			"type": "body",
			"keyword": "/jooyea/images/sns_idea1.jpg"
		}]
	},
	{
		"cms": "iWebSNS",
		"rules": [{
			"type": "body",
			"keyword": "/jooyea/images/snslogo.gif"
		}]
	},
	{
		"cms": "Astaro-Command-Center",
		"rules": [{
			"type": "body",
			"keyword": "/js/_variables_from_backend.js?"
		}]
	},
	{
		"cms": "Astaro-Security-Gateway",
		"rules": [{
			"type": "body",
			"keyword": "/js/_variables_from_backend.js?t="
		}]
	},
	{
		"cms": "CMSTop",
		"rules": [{
			"type": "body",
			"keyword": "/js/cmstop-common.js"
		}]
	},
	{
		"cms": "帝友P2P",
		"rules": [{
			"type": "body",
			"keyword": "/js/diyou.js"
		}]
	},
	{
		"cms": "esoTalk",
		"rules": [{
			"type": "body",
			"keyword": "/js/esotalk.js"
		}]
	},
	{
		"cms": "贷齐乐p2p",
		"rules": [{
			"type": "body",
			"keyword": "/js/jPackageCss/jPackage.css"
		}]
	},
	{
		"cms": "泛微OA-java",
		"rules": [{
			"type": "body",
			"keyword": "/js/jquery/jquery_wev8.js"
		}]
	},
	{
		"cms": "JTBC(CMS)",
		"rules": [{
			"type": "body",
			"keyword": "/js/jtbc.js"
		}]
	},
	{
		"cms": "天融信脆弱性扫描与管理系统",
		"rules": [{
			"type": "body",
			"keyword": "/js/report/horizontalReportPanel.js"
		}]
	},
	{
		"cms": "网动云视讯平台",
		"rules": [{
			"type": "body",
			"keyword": "/js/roomHeight.js"
		}]
	},
	{
		"cms": "rap",
		"rules": [{
			"type": "body",
			"keyword": "/jscripts/rap_util.js"
		}]
	},
	{
		"cms": "rap",
		"rules": [{
			"type": "body",
			"keyword": "/jscripts/rap_util.js"
		}]
	},
	{
		"cms": "金蝶政务GSiS",
		"rules": [{
			"type": "body",
			"keyword": "/kdgs/script/kdgs.js"
		}]
	},
	{
		"cms": "kesionCMS",
		"rules": [{
			"type": "body",
			"keyword": "/ks_inc/common.js"
		}]
	},
	{
		"cms": "Zotonic",
		"rules": [{
			"type": "body",
			"keyword": "/lib/js/apps/zotonic-1.0"
		}]
	},
	{
		"cms": "泛微OA-java",
		"rules": [{
			"type": "body",
			"keyword": "/login/Login.jsp?logintype=1"
		}]
	},
	{
		"cms": "ZKAccess 门禁管理系统",
		"rules": [{
			"type": "body",
			"keyword": "/logoZKAccess_zh-cn.jpg"
		}]
	},
	{
		"cms": "loyaa信息自动采编系统",
		"rules": [{
			"type": "body",
			"keyword": "/Loyaa/common.lib.js"
		}]
	},
	{
		"cms": "Infomaster",
		"rules": [{
			"type": "body",
			"keyword": "/MasterView.css"
		}]
	},
	{
		"cms": "Infomaster",
		"rules": [{
			"type": "body",
			"keyword": "/masterView.js"
		}]
	},
	{
		"cms": "Infomaster",
		"rules": [{
			"type": "body",
			"keyword": "/MasterView/MPLeftNavStyle/PanelBar.MPIFMA.css"
		}]
	},
	{
		"cms": "hikashop",
		"rules": [{
			"type": "body",
			"keyword": "/media/com_hikashop/css/"
		}]
	},
	{
		"cms": "Bulletlink-Newspaper-Template",
		"rules": [{
			"type": "body",
			"keyword": "/ModalPopup/core-modalpopup.css"
		}]
	},
	{
		"cms": "用友erp-nc",
		"rules": [{
			"type": "body",
			"keyword": "/nc/servlet/nc.ui.iufo.login.Index"
		}]
	},
	{
		"cms": "华为 NetOpen",
		"rules": [{
			"type": "body",
			"keyword": "/netopen/theme/css/inFrame.css"
		}]
	},
	{
		"cms": "久其通用财表系统",
		"rules": [{
			"type": "body",
			"keyword": "/netrep/intf"
		}]
	},
	{
		"cms": "久其通用财表系统",
		"rules": [{
			"type": "body",
			"keyword": "/netrep/message2/"
		}]
	},
	{
		"cms": "eagleeyescctv",
		"rules": [{
			"type": "body",
			"keyword": "/nobody/loginDevice.js"
		}]
	},
	{
		"cms": "华天动力OA(OA8000)",
		"rules": [{
			"type": "body",
			"keyword": "/OAapp/WebObjects/OAapp.woa"
		}]
	},
	{
		"cms": "Apache-Wicket",
		"rules": [{
			"type": "body",
			"keyword": "/org.apache.wicket."
		}]
	},
	{
		"cms": "GeoServer",
		"rules": [{
			"type": "body",
			"keyword": "/org.geoserver.web.GeoServerBasePage/"
		}]
	},
	{
		"cms": "ASPilot-Cart",
		"rules": [{
			"type": "body",
			"keyword": "/pilot_css_default.css"
		}]
	},
	{
		"cms": "b2evolution",
		"rules": [{
			"type": "body",
			"keyword": "/powered-by-b2evolution-150t.gif"
		}]
	},
	{
		"cms": "Ultra_Electronics",
		"rules": [{
			"type": "body",
			"keyword": "/preauth/login.cgi"
		}]
	},
	{
		"cms": "Ultra_Electronics",
		"rules": [{
			"type": "body",
			"keyword": "/preauth/style.css"
		}]
	},
	{
		"cms": "悟空CRM",
		"rules": [{
			"type": "body",
			"keyword": "/Public/js/5kcrm.js"
		}]
	},
	{
		"cms": "SCADA PLC",
		"rules": [{
			"type": "body",
			"keyword": "/ralogo.gif"
		}]
	},
	{
		"cms": "xfinity",
		"rules": [{
			"type": "body",
			"keyword": "/reset-meyer-1.0.min.css"
		}]
	},
	{
		"cms": "78oa",
		"rules": [{
			"type": "body",
			"keyword": "/resource/javascript/system/runtime.min.js"
		}]
	},
	{
		"cms": "richmail",
		"rules": [{
			"type": "body",
			"keyword": "/resource/se/lang/se/mail_zh_CN.js"
		}]
	},
	{
		"cms": "iWebShop",
		"rules": [{
			"type": "body",
			"keyword": "/runtime/default/systemjs"
		}]
	},
	{
		"cms": "Siemens_SIMATIC",
		"rules": [{
			"type": "body",
			"keyword": "/S7Web.css"
		}]
	},
	{
		"cms": "蓝凌EIS智慧协同平台",
		"rules": [{
			"type": "body",
			"keyword": "/scripts/jquery.landray.common.js"
		}]
	},
	{
		"cms": "用友致远oa",
		"rules": [{
			"type": "body",
			"keyword": "/seeyon/common/all-min.js"
		}]
	},
	{
		"cms": "OA(a8/seeyon/ufida)",
		"rules": [{
			"type": "body",
			"keyword": "/seeyon/USER-DATA/IMAGES/LOGIN/login.gif"
		}]
	},
	{
		"cms": "用友致远oa",
		"rules": [{
			"type": "body",
			"keyword": "/seeyon/USER-DATA/IMAGES/LOGIN/login.gif"
		}]
	},
	{
		"cms": "IBM-BladeCenter",
		"rules": [{
			"type": "body",
			"keyword": "/shared/ibmbch.png"
		}]
	},
	{
		"cms": "IBM-BladeCenter",
		"rules": [{
			"type": "body",
			"keyword": "/shared/ibmbcs.png"
		}]
	},
	{
		"cms": "AlstraSoft-EPay-Enterprise",
		"rules": [{
			"type": "body",
			"keyword": "/shop.htm?action=view"
		}]
	},
	{
		"cms": "1und1",
		"rules": [{
			"type": "body",
			"keyword": "/shop/catalog/browse?sessid="
		}]
	},
	{
		"cms": "校园卡管理系统",
		"rules": [{
			"type": "body",
			"keyword": "/shouyeziti.css"
		}]
	},
	{
		"cms": "Basilic",
		"rules": [{
			"type": "body",
			"keyword": "/Software/Basilic"
		}]
	},
	{
		"cms": "南方数据",
		"rules": [{
			"type": "body",
			"keyword": "/Southidcj2f.Js"
		}]
	},
	{
		"cms": "南方数据",
		"rules": [{
			"type": "body",
			"keyword": "/SouthidcKeFu.js"
		}]
	},
	{
		"cms": "金山KingGate",
		"rules": [{
			"type": "body",
			"keyword": "/src/system/login.php"
		}]
	},
	{
		"cms": "CDR-Stats",
		"rules": [{
			"type": "body",
			"keyword": "/static/cdr-stats/js/jquery"
		}]
	},
	{
		"cms": "睿博士云办公系统",
		"rules": [{
			"type": "body",
			"keyword": "/studentSign/toLogin.di"
		}]
	},
	{
		"cms": "逐浪zoomla",
		"rules": [{
			"type": "body",
			"keyword": "/style/images/win8_symbol_140x140.png"
		}]
	},
	{
		"cms": "用友ufida",
		"rules": [{
			"type": "body",
			"keyword": "/System/Login/Login.asp?AppID="
		}]
	},
	{
		"cms": "科信邮件系统",
		"rules": [{
			"type": "body",
			"keyword": "/systemfunction.pack.js"
		}]
	},
	{
		"cms": "DedeCMS",
		"rules": [{
			"type": "body",
			"keyword": "/templets/default/style/dedecms.css"
		}]
	},
	{
		"cms": "微门户",
		"rules": [{
			"type": "body",
			"keyword": "/tpl/Home/weimeng/common/css/"
		}]
	},
	{
		"cms": "睿博士云办公系统",
		"rules": [{
			"type": "body",
			"keyword": "/user/toUpdatePasswordPage.di"
		}]
	},
	{
		"cms": "easypanel",
		"rules": [{
			"type": "body",
			"keyword": "/vhost/view/default/style/login.css"
		}]
	},
	{
		"cms": "佳能网络摄像头(Canon Network Cameras)",
		"rules": [{
			"type": "body",
			"keyword": "/viewer/live/en/live.html"
		}]
	},
	{
		"cms": "trs_wcm",
		"rules": [{
			"type": "body",
			"keyword": "/wcm\" target=\"_blank\">管理"
		}]
	},
	{
		"cms": "trs_wcm",
		"rules": [{
			"type": "body",
			"keyword": "/wcm\" target=\"_blank\">网站管理"
		}]
	},
	{
		"cms": "trs_wcm",
		"rules": [{
			"type": "body",
			"keyword": "/wcm/app/js"
		}]
	},
	{
		"cms": "Avaya-Aura-Utility-Server",
		"rules": [{
			"type": "body",
			"keyword": "/webhelp/Base/Utility_toc.htm"
		}]
	},
	{
		"cms": "we7",
		"rules": [{
			"type": "body",
			"keyword": "/Widgets/WidgetCollection/"
		}]
	},
	{
		"cms": "MediaWiki",
		"rules": [{
			"type": "body",
			"keyword": "/wiki/images/6/64/Favicon.ico"
		}]
	},
	{
		"cms": "mirapoint",
		"rules": [{
			"type": "body",
			"keyword": "/wm/mail/login.html"
		}]
	},
	{
		"cms": "MDaemon",
		"rules": [{
			"type": "body",
			"keyword": "/WorldClient.dll?View=Main"
		}]
	},
	{
		"cms": "WordPress-php",
		"rules": [{
			"type": "body",
			"keyword": "/wp-login.php?"
		}]
	},
	{
		"cms": "苏亚星校园管理系统",
		"rules": [{
			"type": "body",
			"keyword": "/ws2004/Public/"
		}]
	},
	{
		"cms": "FortiGuard",
		"rules": [{
			"type": "body",
			"keyword": "/XX/YY/ZZ/CI/MGPGHGPGPFGHCDPFGGOGFGEH"
		}]
	},
	{
		"cms": "一米OA",
		"rules": [{
			"type": "body",
			"keyword": "/yimioa.apk"
		}]
	},
	{
		"cms": "用友致远oa",
		"rules": [{
			"type": "body",
			"keyword": "/yyoa/"
		}]
	},
	{
		"cms": "zenoss",
		"rules": [{
			"type": "body",
			"keyword": "/zport/dmd/"
		}]
	},
	{
		"cms": "WP Plugin All-in-one-SEO-Pack",
		"rules": [{
			"type": "body",
			"keyword": "<!-- /all in one seo pack -->"
		}]
	},
	{
		"cms": "EDK",
		"rules": [{
			"type": "body",
			"keyword": "<!-- /killlistable.tpl -->"
		}]
	},
	{
		"cms": "Aurion",
		"rules": [{
			"type": "body",
			"keyword": "<!-- Aurion Teal will be used as the login-time default"
		}]
	},
	{
		"cms": "BASE",
		"rules": [{
			"type": "body",
			"keyword": "<!-- Basic Analysis and Security Engine (BASE) -->"
		}]
	},
	{
		"cms": "CaupoShop-Classic",
		"rules": [{
			"type": "body",
			"keyword": "<!-- CaupoShop Classic"
		}]
	},
	{
		"cms": "ClipBucket",
		"rules": [{
			"type": "body",
			"keyword": "<!-- ClipBucket"
		}]
	},
	{
		"cms": "ClipBucket",
		"rules": [{
			"type": "body",
			"keyword": "<!-- Forged by ClipBucket"
		}]
	},
	{
		"cms": "coWiki",
		"rules": [{
			"type": "body",
			"keyword": "<!-- Generated by coWiki"
		}]
	},
	{
		"cms": "bbPress",
		"rules": [{
			"type": "body",
			"keyword": "<!-- If you like showing off the fact that your server rocks -->"
		}]
	},
	{
		"cms": "Intraxxion-CMS",
		"rules": [{
			"type": "body",
			"keyword": "<!-- site built by Intraxxion"
		}]
	},
	{
		"cms": "CA-SiteMinder",
		"rules": [{
			"type": "body",
			"keyword": "<!-- SiteMinder Encoding"
		}]
	},
	{
		"cms": "Escenic",
		"rules": [{
			"type": "body",
			"keyword": "<!-- Start Escenic Analysis Engine client script -->"
		}]
	},
	{
		"cms": "BM-Classifieds",
		"rules": [{
			"type": "body",
			"keyword": "<!-- START HEADER TABLE - HOLDS GRAPHIC AND SITE NAME -->"
		}]
	},
	{
		"cms": "CGI:IRC",
		"rules": [{
			"type": "body",
			"keyword": "<!-- This is part of CGI:IRC"
		}]
	},
	{
		"cms": "ClipShare",
		"rules": [{
			"type": "body",
			"keyword": "<!--!!!!!!!!!!!!!!!!!!!!!!!!! Processing SCRIPT"
		}]
	},
	{
		"cms": "DZCP",
		"rules": [{
			"type": "body",
			"keyword": "<!--[ DZCP"
		}]
	},
	{
		"cms": "Coppermine",
		"rules": [{
			"type": "body",
			"keyword": "<!--Coppermine Photo Gallery"
		}]
	},
	{
		"cms": "EazyCMS",
		"rules": [{
			"type": "body",
			"keyword": "<a class=\"actionlink\" href=\"http://www.eazyCMS.com"
		}]
	},
	{
		"cms": "CMSTop",
		"rules": [{
			"type": "body",
			"keyword": "<a class=\"poweredby\" href=\"http://www.cmstop.com\""
		}]
	},
	{
		"cms": "eFront",
		"rules": [{
			"type": "body",
			"keyword": "<a href = \"http://www.efrontlearning.net"
		}]
	},
	{
		"cms": "mongodb",
		"rules": [{
			"type": "body",
			"keyword": "<a href=\"/_replSet\">Replica set status</a></p>"
		}]
	},
	{
		"cms": "Fossil",
		"rules": [{
			"type": "body",
			"keyword": "<a href=\"http://fossil-scm.org"
		}]
	},
	{
		"cms": "Bomgar",
		"rules": [{
			"type": "body",
			"keyword": "<a href=\"http://www.bomgar.com/products\" class=\"inverse"
		}]
	},
	{
		"cms": "CaupoShop-Classic",
		"rules": [{
			"type": "body",
			"keyword": "<a href=\"http://www.caupo.net\" target=\"_blank\">CaupoNet"
		}]
	},
	{
		"cms": "DSpace",
		"rules": [{
			"type": "body",
			"keyword": "<a href=\"http://www.dspace.org\">DSpace Software"
		}]
	},
	{
		"cms": "eTicket",
		"rules": [{
			"type": "body",
			"keyword": "<a href=\"http://www.eticketsupport.com\" target=\"_blank\">"
		}]
	},
	{
		"cms": "FileVista",
		"rules": [{
			"type": "body",
			"keyword": "<a href=\"http://www.gleamtech.com/products/filevista/web-file-manager"
		}]
	},
	{
		"cms": "gCards",
		"rules": [{
			"type": "body",
			"keyword": "<a href=\"http://www.gregphoto.net/gcards/index.php"
		}]
	},
	{
		"cms": "GridSite",
		"rules": [{
			"type": "body",
			"keyword": "<a href=\"http://www.gridsite.org/\">GridSite"
		}]
	},
	{
		"cms": "CGIProxy",
		"rules": [{
			"type": "body",
			"keyword": "<a href=\"http://www.jmarshall.com/tools/cgiproxy/"
		}]
	},
	{
		"cms": "AlstraSoft-AskMe",
		"rules": [{
			"type": "body",
			"keyword": "<a href=\"pass_recover.php\">"
		}]
	},
	{
		"cms": "Citrix-XenServer",
		"rules": [{
			"type": "body",
			"keyword": "<a href=\"XenCenterSetup.exe\">XenCenter installer</a>"
		}]
	},
	{
		"cms": "CafeEngine",
		"rules": [{
			"type": "body",
			"keyword": "<a href=http://cafeengine.com>CafeEngine.com"
		}]
	},
	{
		"cms": "Evo-Cam",
		"rules": [{
			"type": "body",
			"keyword": "<applet archive=\"evocam.jar"
		}]
	},
	{
		"cms": "Carrier-CCNWeb",
		"rules": [{
			"type": "body",
			"keyword": "<APPLET CODE=\"JLogin.class\" ARCHIVE=\"JLogin.jar"
		}]
	},
	{
		"cms": "U-Mail",
		"rules": [{
			"type": "body",
			"keyword": "<BODY LINK=\"White\" VLINK=\"White\" ALINK=\"White\">"
		}]
	},
	{
		"cms": "php云",
		"rules": [{
			"type": "body",
			"keyword": "<div class=\"index_link_list_name\">"
		}]
	},
	{
		"cms": "BugTracker.NET",
		"rules": [{
			"type": "body",
			"keyword": "<div class=logo>BugTracker.NET"
		}]
	},
	{
		"cms": "BackBee",
		"rules": [{
			"type": "body",
			"keyword": "<div id=\"bb5-site-wrapper\">"
		}]
	},
	{
		"cms": "DokuWiki",
		"rules": [{
			"type": "body",
			"keyword": "<div id=\"dokuwiki"
		}]
	},
	{
		"cms": "ICEshop",
		"rules": [{
			"type": "body",
			"keyword": "<div id=\"iceshop\">"
		}]
	},
	{
		"cms": "DnP-Firewall",
		"rules": [{
			"type": "body",
			"keyword": "<form name=dnp_firewall"
		}]
	},
	{
		"cms": "awstats_admin",
		"rules": [{
			"type": "body",
			"keyword": "<frame name=\"mainleft\" src=\"awstats.pl?config="
		}]
	},
	{
		"cms": "Brother-Printer",
		"rules": [{
			"type": "body",
			"keyword": "<FRAME SRC=\"/printer/inc_head.html"
		}]
	},
	{
		"cms": "BEA-WebLogic-Server",
		"rules": [{
			"type": "body",
			"keyword": "<h1>BEA WebLogic Server"
		}]
	},
	{
		"cms": "STAR CMS",
		"rules": [{
			"type": "body",
			"keyword": "<img alt=\"STAR CMS"
		}]
	},
	{
		"cms": "Brother-Printer",
		"rules": [{
			"type": "body",
			"keyword": "<IMG src=\"/common/image/HL4040CN"
		}]
	},
	{
		"cms": "File-Upload-Manager",
		"rules": [{
			"type": "body",
			"keyword": "<IMG SRC=\"/images/header.jpg\" ALT=\"File Upload Manager\">"
		}]
	},
	{
		"cms": "中国期刊先知网",
		"rules": [{
			"type": "body",
			"keyword": "<img src=\"images/logoknow.png\""
		}]
	},
	{
		"cms": "Atmail-WebMail",
		"rules": [{
			"type": "body",
			"keyword": "<input id=\"Mailserverinput"
		}]
	},
	{
		"cms": "Apabi数字资源平台",
		"rules": [{
			"type": "body",
			"keyword": "<link href=\"HTTP://apabi"
		}]
	},
	{
		"cms": "通达OA",
		"rules": [{
			"type": "body",
			"keyword": "<link rel=\"shortcut icon\" href=\"/images/tongda.ico\" />"
		}]
	},
	{
		"cms": "VOS3000",
		"rules": [{
			"type": "body",
			"keyword": "<meta name=\"description\" content=\"VOS3000"
		}]
	},
	{
		"cms": "destoon",
		"rules": [{
			"type": "body",
			"keyword": "<meta name=\"generator\" content=\"Destoon"
		}]
	},
	{
		"cms": "VOS3000",
		"rules": [{
			"type": "body",
			"keyword": "<meta name=\"keywords\" content=\"VOS3000"
		}]
	},
	{
		"cms": "久其通用财表系统",
		"rules": [{
			"type": "body",
			"keyword": "<nobr>北京久其软件股份有限公司"
		}]
	},
	{
		"cms": "activeCollab",
		"rules": [{
			"type": "body",
			"keyword": "<p id=\"powered_by\"><a href=\"http://www.activecollab.com/\""
		}]
	},
	{
		"cms": "Foxycart",
		"rules": [{
			"type": "body",
			"keyword": "<script src=\"//cdn.foxycart.com"
		}]
	},
	{
		"cms": "CGI:IRC",
		"rules": [{
			"type": "body",
			"keyword": "<small id=\"ietest\"><a href=\"http://cgiirc.org/"
		}]
	},
	{
		"cms": "CitusCMS",
		"rules": [{
			"type": "body",
			"keyword": "<strong>CitusCMS</strong>"
		}]
	},
	{
		"cms": "HostBill",
		"rules": [{
			"type": "body",
			"keyword": "<strong>HostBill"
		}]
	},
	{
		"cms": "宝塔面板-python",
		"rules": [{
			"type": "body",
			"keyword": "<title>安全入口校验失败</title>"
		}]
	},
	{
		"cms": "Amiro-CMS",
		"rules": [{
			"type": "body",
			"keyword": "-= Amiro.CMS (c) =-"
		}]
	},
	{
		"cms": "Buddy-Zone",
		"rules": [{
			"type": "body",
			"keyword": ">Buddy Zone</a>"
		}]
	},
	{
		"cms": "taocms",
		"rules": [{
			"type": "body",
			"keyword": ">taoCMS<"
		}]
	},
	{
		"cms": "TPshop",
		"rules": [{
			"type": "body",
			"keyword": ">TPshop开源商城<"
		}]
	},
	{
		"cms": "trs_wcm",
		"rules": [{
			"type": "body",
			"keyword": "0;URL=/wcm"
		}]
	},
	{
		"cms": "农友政务系统",
		"rules": [{
			"type": "body",
			"keyword": "1207044504"
		}]
	},
	{
		"cms": "DVR-WebClient",
		"rules": [{
			"type": "body",
			"keyword": "259F9FDF-97EA-4C59-B957-5160CAB6884E"
		}]
	},
	{
		"cms": "360企业版",
		"rules": [{
			"type": "body",
			"keyword": "360EntInst"
		}]
	},
	{
		"cms": "360webfacil_360WebManager",
		"rules": [{
			"type": "body",
			"keyword": "360WebManager Software"
		}]
	},
	{
		"cms": "huawei_auth_server",
		"rules": [{
			"type": "body",
			"keyword": "75718C9A-F029-11d1-A1AC-00C04FB6C223"
		}]
	},
	{
		"cms": "PaloAlto_Firewall",
		"rules": [{
			"type": "body",
			"keyword": "Access to the web page you were trying to visit has been blocked in accordance with company policy"
		}]
	},
	{
		"cms": "ACTi",
		"rules": [{
			"type": "body",
			"keyword": "ACTi Corporation All Rights Reserved"
		}]
	},
	{
		"cms": "PHPOA",
		"rules": [{
			"type": "body",
			"keyword": "admin_img/msg_bg.png"
		}]
	},
	{
		"cms": "协众OA",
		"rules": [{
			"type": "body",
			"keyword": "admin@cnoa.cn"
		}]
	},
	{
		"cms": "HP-OfficeJet-Printer",
		"rules": [{
			"type": "body",
			"keyword": "align=\"center\">HP Officejet"
		}]
	},
	{
		"cms": "福富安全基线管理",
		"rules": [{
			"type": "body",
			"keyword": "align=\"center\">福富软件"
		}]
	},
	{
		"cms": "IBM-BladeCenter",
		"rules": [{
			"type": "body",
			"keyword": "alt=\"IBM BladeCenter"
		}]
	},
	{
		"cms": "ionCube-Loader",
		"rules": [{
			"type": "body",
			"keyword": "alt=\"ionCube logo"
		}]
	},
	{
		"cms": "Bomgar",
		"rules": [{
			"type": "body",
			"keyword": "alt=\"Remote Support by BOMGAR"
		}]
	},
	{
		"cms": "汉码软件",
		"rules": [{
			"type": "body",
			"keyword": "alt=\"汉码软件LOGO"
		}]
	},
	{
		"cms": "AlumniServer",
		"rules": [{
			"type": "body",
			"keyword": "AlumniServerProject.php"
		}]
	},
	{
		"cms": "AM4SS",
		"rules": [{
			"type": "body",
			"keyword": "am4ss.css"
		}]
	},
	{
		"cms": "Array_Networks_VPN",
		"rules": [{
			"type": "body",
			"keyword": "an_util.js"
		}]
	},
	{
		"cms": "AiCart",
		"rules": [{
			"type": "body",
			"keyword": "APP_authenticate"
		}]
	},
	{
		"cms": "Solr",
		"rules": [{
			"type": "body",
			"keyword": "app_config.solr_path"
		}]
	},
	{
		"cms": "擎天电子政务",
		"rules": [{
			"type": "body",
			"keyword": "App_Themes/1/Style.css"
		}]
	},
	{
		"cms": "AppServ",
		"rules": [{
			"type": "body",
			"keyword": "appserv/softicon.gif"
		}]
	},
	{
		"cms": "ASPThai_Net-Webboard",
		"rules": [{
			"type": "body",
			"keyword": "ASPThai.Net Webboard"
		}]
	},
	{
		"cms": "微普外卖点餐系统",
		"rules": [{
			"type": "body",
			"keyword": "Author\" content=\"微普外卖点餐系统"
		}]
	},
	{
		"cms": "Munin",
		"rules": [{
			"type": "body",
			"keyword": "Auto-generated by Munin"
		}]
	},
	{
		"cms": "wecenter",
		"rules": [{
			"type": "body",
			"keyword": "aw_template.js"
		}]
	},
	{
		"cms": "awstats",
		"rules": [{
			"type": "body",
			"keyword": "awstats.pl?config="
		}]
	},
	{
		"cms": "axis2-web",
		"rules": [{
			"type": "body",
			"keyword": "axis2-web/css/axis-style.css"
		}]
	},
	{
		"cms": "蓝盾BDWebGuard",
		"rules": [{
			"type": "body",
			"keyword": "BACKGROUND: url(images/loginbg.jpg) #e5f1fc"
		}]
	},
	{
		"cms": "bluecms",
		"rules": [{
			"type": "body",
			"keyword": "bcms_plugin"
		}]
	},
	{
		"cms": "BigDump",
		"rules": [{
			"type": "body",
			"keyword": "BigDump: Staggered MySQL Dump Importer"
		}]
	},
	{
		"cms": "北京阳光环球建站系统",
		"rules": [{
			"type": "body",
			"keyword": "bigSortProduct.asp?bigid"
		}]
	},
	{
		"cms": "bit-service",
		"rules": [{
			"type": "body",
			"keyword": "bit-xxzs"
		}]
	},
	{
		"cms": "mantis",
		"rules": [{
			"type": "body",
			"keyword": "browser_search_plugin.php?type=id"
		}]
	},
	{
		"cms": "ASProxy",
		"rules": [{
			"type": "body",
			"keyword": "btnASProxyDisplayButton"
		}]
	},
	{
		"cms": "Dolphin",
		"rules": [{
			"type": "body",
			"keyword": "bx_css_async"
		}]
	},
	{
		"cms": "Hiki",
		"rules": [{
			"type": "body",
			"keyword": "by <a href=\"http://hikiwiki.org/"
		}]
	},
	{
		"cms": "glFusion",
		"rules": [{
			"type": "body",
			"keyword": "by <a href=\"http://www.glfusion.org/"
		}]
	},
	{
		"cms": "Imageview",
		"rules": [{
			"type": "body",
			"keyword": "By Jorge Schrauwen"
		}]
	},
	{
		"cms": "OpenCart",
		"rules": [{
			"type": "body",
			"keyword": "catalog/view/theme"
		}]
	},
	{
		"cms": "亿赛通DLP",
		"rules": [{
			"type": "body",
			"keyword": "CDGServer3"
		}]
	},
	{
		"cms": "Cisco-IP-Phone",
		"rules": [{
			"type": "body",
			"keyword": "Cisco Unified Wireless IP Phone"
		}]
	},
	{
		"cms": "Citrix-XenServer",
		"rules": [{
			"type": "body",
			"keyword": "Citrix Systems, Inc. XenServer"
		}]
	},
	{
		"cms": "Car-Portal",
		"rules": [{
			"type": "body",
			"keyword": "class=\"bodyfontwhite\"><strong>&nbsp;Car Script"
		}]
	},
	{
		"cms": "GuppY",
		"rules": [{
			"type": "body",
			"keyword": "class=\"copyright\" href=\"http://www.freeguppy.org/"
		}]
	},
	{
		"cms": "GeoServer",
		"rules": [{
			"type": "body",
			"keyword": "class=\"geoserver lebeg"
		}]
	},
	{
		"cms": "DV-Cart",
		"rules": [{
			"type": "body",
			"keyword": "class=\"KT_tngtable"
		}]
	},
	{
		"cms": "BugFree",
		"rules": [{
			"type": "body",
			"keyword": "class=\"loginBgImage\" alt=\"BugFree"
		}]
	},
	{
		"cms": "DaDaBIK",
		"rules": [{
			"type": "body",
			"keyword": "class=\"powered_by_dadabik"
		}]
	},
	{
		"cms": "Interspire-Shopping-Cart",
		"rules": [{
			"type": "body",
			"keyword": "class=\"PoweredBy\">Interspire Shopping Cart"
		}]
	},
	{
		"cms": "xheditor",
		"rules": [{
			"type": "body",
			"keyword": "class=\"xheditor"
		}]
	},
	{
		"cms": "dasannetworks",
		"rules": [{
			"type": "body",
			"keyword": "clear_cookie(\"login\");"
		}]
	},
	{
		"cms": "Comcast_Business",
		"rules": [{
			"type": "body",
			"keyword": "cmn/css/common-min.css"
		}]
	},
	{
		"cms": "CMSTop",
		"rules": [{
			"type": "body",
			"keyword": "cmstop-list-text.css"
		}]
	},
	{
		"cms": "IBM-Cognos",
		"rules": [{
			"type": "body",
			"keyword": "Cognos &#26159; International Business Machines Corp"
		}]
	},
	{
		"cms": "i@Report",
		"rules": [{
			"type": "body",
			"keyword": "com.sanlink.server.Login"
		}]
	},
	{
		"cms": "Comcast_Business_Gateway",
		"rules": [{
			"type": "body",
			"keyword": "Comcast Business Gateway"
		}]
	},
	{
		"cms": "MRTG",
		"rules": [{
			"type": "body",
			"keyword": "Command line is easier to read using \"View Page Properties\" of your browser"
		}]
	},
	{
		"cms": "MRTG",
		"rules": [{
			"type": "body",
			"keyword": "commandline was: indexmaker"
		}]
	},
	{
		"cms": "08cms",
		"rules": [{
			"type": "body",
			"keyword": "content=\"08CMS"
		}]
	},
	{
		"cms": "1024 CMS",
		"rules": [{
			"type": "body",
			"keyword": "content=\"1024 CMS"
		}]
	},
	{
		"cms": "171cms",
		"rules": [{
			"type": "body",
			"keyword": "content=\"171cms"
		}]
	},
	{
		"cms": "任我行电商",
		"rules": [{
			"type": "body",
			"keyword": "content=\"366EC"
		}]
	},
	{
		"cms": "74cms",
		"rules": [{
			"type": "body",
			"keyword": "content=\"74cms.com"
		}]
	},
	{
		"cms": "AChecker Web accessibility evaluation tool",
		"rules": [{
			"type": "body",
			"keyword": "content=\"AChecker is a Web accessibility"
		}]
	},
	{
		"cms": "Allomani",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Allomani"
		}]
	},
	{
		"cms": "AlumniServer",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Alumni"
		}]
	},
	{
		"cms": "Apache-Forrest",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Apache Forrest"
		}]
	},
	{
		"cms": "ASP-Nuke",
		"rules": [{
			"type": "body",
			"keyword": "CONTENT=\"ASP-Nuke"
		}]
	},
	{
		"cms": "ASPCMS",
		"rules": [{
			"type": "body",
			"keyword": "content=\"ASPCMS"
		}]
	},
	{
		"cms": "ASP-Nuke",
		"rules": [{
			"type": "body",
			"keyword": "content=\"ASPNUKE"
		}]
	},
	{
		"cms": "Auto-CMS",
		"rules": [{
			"type": "body",
			"keyword": "content=\"AutoCMS"
		}]
	},
	{
		"cms": "Axous",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Axous"
		}]
	},
	{
		"cms": "b2bbuilder",
		"rules": [{
			"type": "body",
			"keyword": "content=\"B2Bbuilder"
		}]
	},
	{
		"cms": "b2evolution",
		"rules": [{
			"type": "body",
			"keyword": "content=\"b2evolution"
		}]
	},
	{
		"cms": "八哥CMS",
		"rules": [{
			"type": "body",
			"keyword": "content=\"BageCMS"
		}]
	},
	{
		"cms": "baocms",
		"rules": [{
			"type": "body",
			"keyword": "content=\"BAOCMS"
		}]
	},
	{
		"cms": "Belkin-Modem",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Belkin"
		}]
	},
	{
		"cms": "BIGACE",
		"rules": [{
			"type": "body",
			"keyword": "content=\"BIGACE"
		}]
	},
	{
		"cms": "bitweaver",
		"rules": [{
			"type": "body",
			"keyword": "content=\"bitweaver"
		}]
	},
	{
		"cms": "BloofoxCMS",
		"rules": [{
			"type": "body",
			"keyword": "content=\"bloofoxCMS"
		}]
	},
	{
		"cms": "SiteEngine",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Boka SiteEngine"
		}]
	},
	{
		"cms": "IQeye-Netcam",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Brian Lau, IQinVision"
		}]
	},
	{
		"cms": "BrowserCMS",
		"rules": [{
			"type": "body",
			"keyword": "content=\"BrowserCMS"
		}]
	},
	{
		"cms": "CameraLife",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Camera Life"
		}]
	},
	{
		"cms": "Campsite",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Campsite"
		}]
	},
	{
		"cms": "CitusCMS",
		"rules": [{
			"type": "body",
			"keyword": "content=\"CitusCMS"
		}]
	},
	{
		"cms": "ClanSphere",
		"rules": [{
			"type": "body",
			"keyword": "content=\"ClanSphere"
		}]
	},
	{
		"cms": "ClipBucket",
		"rules": [{
			"type": "body",
			"keyword": "content=\"ClipBucket"
		}]
	},
	{
		"cms": "CMScontrol",
		"rules": [{
			"type": "body",
			"keyword": "content=\"CMScontrol"
		}]
	},
	{
		"cms": "CMSimple",
		"rules": [{
			"type": "body",
			"keyword": "content=\"CMSimple"
		}]
	},
	{
		"cms": "CommonSpot",
		"rules": [{
			"type": "body",
			"keyword": "content=\"CommonSpot"
		}]
	},
	{
		"cms": "ContentXXL",
		"rules": [{
			"type": "body",
			"keyword": "content=\"contentXXL"
		}]
	},
	{
		"cms": "Contrexx-CMS",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Contrexx"
		}]
	},
	{
		"cms": "南方数据",
		"rules": [{
			"type": "body",
			"keyword": "CONTENT=\"Copyright 2003-2015 - Southidc.net"
		}]
	},
	{
		"cms": "coWiki",
		"rules": [{
			"type": "body",
			"keyword": "content=\"coWiki"
		}]
	},
	{
		"cms": "Custom-CMS",
		"rules": [{
			"type": "body",
			"keyword": "content=\"CustomCMS"
		}]
	},
	{
		"cms": "Cyn_in",
		"rules": [{
			"type": "body",
			"keyword": "content=\"cyn.in"
		}]
	},
	{
		"cms": "DaDaBIK",
		"rules": [{
			"type": "body",
			"keyword": "content=\"DaDaBIK"
		}]
	},
	{
		"cms": "大米CMS",
		"rules": [{
			"type": "body",
			"keyword": "content=\"damicms"
		}]
	},
	{
		"cms": "dbshop",
		"rules": [{
			"type": "body",
			"keyword": "content=\"dbshop"
		}]
	},
	{
		"cms": "DirCMS",
		"rules": [{
			"type": "body",
			"keyword": "content=\"DirCMS"
		}]
	},
	{
		"cms": "discuz",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Discuz"
		}]
	},
	{
		"cms": "Dokeos",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Dokeos"
		}]
	},
	{
		"cms": "DokuWiki",
		"rules": [{
			"type": "body",
			"keyword": "content=\"DokuWiki"
		}]
	},
	{
		"cms": "DORG",
		"rules": [{
			"type": "body",
			"keyword": "CONTENT=\"DORG"
		}]
	},
	{
		"cms": "DotA-OpenStats",
		"rules": [{
			"type": "body",
			"keyword": "content=\"dota OpenStats"
		}]
	},
	{
		"cms": "DSpace",
		"rules": [{
			"type": "body",
			"keyword": "content=\"DSpace"
		}]
	},
	{
		"cms": "dswjcms",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Dswjcms"
		}]
	},
	{
		"cms": "DT-Centrepiece",
		"rules": [{
			"type": "body",
			"keyword": "content=\"DT Centrepiece"
		}]
	},
	{
		"cms": "eadmin",
		"rules": [{
			"type": "body",
			"keyword": "content=\"eAdmin"
		}]
	},
	{
		"cms": "easyLink-Web-Solutions",
		"rules": [{
			"type": "body",
			"keyword": "content=\"easyLink"
		}]
	},
	{
		"cms": "Ecomat-CMS",
		"rules": [{
			"type": "body",
			"keyword": "content=\"ECOMAT CMS"
		}]
	},
	{
		"cms": "EDIMAX",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Edimax"
		}]
	},
	{
		"cms": "Edito-CMS",
		"rules": [{
			"type": "body",
			"keyword": "content=\"edito"
		}]
	},
	{
		"cms": "eGroupWare",
		"rules": [{
			"type": "body",
			"keyword": "content=\"eGroupWare"
		}]
	},
	{
		"cms": "eLitius",
		"rules": [{
			"type": "body",
			"keyword": "content=\"eLitius"
		}]
	},
	{
		"cms": "Elxis-CMS",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Elxis"
		}]
	},
	{
		"cms": "EPiServer",
		"rules": [{
			"type": "body",
			"keyword": "content=\"EPiServer"
		}]
	},
	{
		"cms": "ANECMS",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Erwin Aligam - ealigam@gmail.com"
		}]
	},
	{
		"cms": "Escenic",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Escenic"
		}]
	},
	{
		"cms": "Contentteller-CMS",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Esselbach Contentteller CMS"
		}]
	},
	{
		"cms": "E-Xoopport",
		"rules": [{
			"type": "body",
			"keyword": "content=\"E-Xoopport"
		}]
	},
	{
		"cms": "eSyndiCat",
		"rules": [{
			"type": "body",
			"keyword": "content=\"eSyndiCat"
		}]
	},
	{
		"cms": "Exponent-CMS",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Exponent Content Management System"
		}]
	},
	{
		"cms": "Fastpublish-CMS",
		"rules": [{
			"type": "body",
			"keyword": "content=\"fastpublish"
		}]
	},
	{
		"cms": "fengcms",
		"rules": [{
			"type": "body",
			"keyword": "content=\"FengCms"
		}]
	},
	{
		"cms": "FluentNET",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Fluent"
		}]
	},
	{
		"cms": "Fluid-Dynamics-Search-Engine",
		"rules": [{
			"type": "body",
			"keyword": "content=\"fluid dynamics"
		}]
	},
	{
		"cms": "KaiBB",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Forum powered by KaiBB"
		}]
	},
	{
		"cms": "FoxPHP",
		"rules": [{
			"type": "body",
			"keyword": "content=\"FoxPHP"
		}]
	},
	{
		"cms": "Gallarific",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Gallarific"
		}]
	},
	{
		"cms": "GetSimple",
		"rules": [{
			"type": "body",
			"keyword": "content=\"GetSimple"
		}]
	},
	{
		"cms": "Glossword",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Glossword"
		}]
	},
	{
		"cms": "GuppY",
		"rules": [{
			"type": "body",
			"keyword": "content=\"GuppY"
		}]
	},
	{
		"cms": "Hiki",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Hiki"
		}]
	},
	{
		"cms": "Hotaru-CMS",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Hotaru"
		}]
	},
	{
		"cms": "Ananyoo-CMS",
		"rules": [{
			"type": "body",
			"keyword": "content=\"http://www.ananyoo.com"
		}]
	},
	{
		"cms": "Hycus-CMS",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Hycus"
		}]
	},
	{
		"cms": "Ikonboard",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Ikonboard"
		}]
	},
	{
		"cms": "Imageview",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Imageview"
		}]
	},
	{
		"cms": "IMGCms",
		"rules": [{
			"type": "body",
			"keyword": "content=\"IMGCMS"
		}]
	},
	{
		"cms": "ImpressPages-CMS",
		"rules": [{
			"type": "body",
			"keyword": "content=\"ImpressPages CMS"
		}]
	},
	{
		"cms": "Informatics-CMS",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Informatics"
		}]
	},
	{
		"cms": "InterRed",
		"rules": [{
			"type": "body",
			"keyword": "content=\"InterRed"
		}]
	},
	{
		"cms": "Interspire-Shopping-Cart",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Interspire Shopping Cart"
		}]
	},
	{
		"cms": "Intraxxion-CMS",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Intraxxion"
		}]
	},
	{
		"cms": "Jamroom",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Jamroom"
		}]
	},
	{
		"cms": "javashop",
		"rules": [{
			"type": "body",
			"keyword": "content=\"JavaShop"
		}]
	},
	{
		"cms": "Jcow",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Jcow"
		}]
	},
	{
		"cms": "jieqi cms",
		"rules": [{
			"type": "body",
			"keyword": "content=\"jieqi cms"
		}]
	},
	{
		"cms": "jobberBase",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Jobberbase"
		}]
	},
	{
		"cms": "Joomla",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Joomla"
		}]
	},
	{
		"cms": "JTBC(CMS)",
		"rules": [{
			"type": "body",
			"keyword": "content=\"JTBC"
		}]
	},
	{
		"cms": "Kajona",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Kajona"
		}]
	},
	{
		"cms": "Kandidat-CMS",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Kandidat-CMS"
		}]
	},
	{
		"cms": "kingcms",
		"rules": [{
			"type": "body",
			"keyword": "content=\"KingCMS"
		}]
	},
	{
		"cms": "lepton-cms",
		"rules": [{
			"type": "body",
			"keyword": "content=\"LEPTON-CMS"
		}]
	},
	{
		"cms": "MallBuilder",
		"rules": [{
			"type": "body",
			"keyword": "content=\"MallBuilder"
		}]
	},
	{
		"cms": "MetInfo",
		"rules": [{
			"type": "body",
			"keyword": "content=\"MetInfo"
		}]
	},
	{
		"cms": "MoMoCMS",
		"rules": [{
			"type": "body",
			"keyword": "content=\"MoMoCMS"
		}]
	},
	{
		"cms": "MvMmall",
		"rules": [{
			"type": "body",
			"keyword": "content=\"MvMmall"
		}]
	},
	{
		"cms": "mymps",
		"rules": [{
			"type": "body",
			"keyword": "content=\"mymps"
		}]
	},
	{
		"cms": "牛逼cms",
		"rules": [{
			"type": "body",
			"keyword": "content=\"niubicms"
		}]
	},
	{
		"cms": "niucms",
		"rules": [{
			"type": "body",
			"keyword": "content=\"NIUCMS"
		}]
	},
	{
		"cms": "opencms",
		"rules": [{
			"type": "body",
			"keyword": "content=\"OpenCms"
		}]
	},
	{
		"cms": "OpenSNS",
		"rules": [{
			"type": "body",
			"keyword": "content=\"OpenSNS"
		}]
	},
	{
		"cms": "DotA-OpenStats",
		"rules": [{
			"type": "body",
			"keyword": "content=\"openstats.iz.rs"
		}]
	},
	{
		"cms": "ourphp",
		"rules": [{
			"type": "body",
			"keyword": "content=\"OURPHP"
		}]
	},
	{
		"cms": "PageAdmin",
		"rules": [{
			"type": "body",
			"keyword": "content=\"PageAdmin CMS\""
		}]
	},
	{
		"cms": "PhpCMS",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Phpcms"
		}]
	},
	{
		"cms": "phpdisk",
		"rules": [{
			"type": "body",
			"keyword": "content=\"PHPDisk"
		}]
	},
	{
		"cms": "phpems考试系统",
		"rules": [{
			"type": "body",
			"keyword": "content=\"PHPEMS"
		}]
	},
	{
		"cms": "PHPMyWind",
		"rules": [{
			"type": "body",
			"keyword": "content=\"PHPMyWind"
		}]
	},
	{
		"cms": "phpok",
		"rules": [{
			"type": "body",
			"keyword": "content=\"phpok"
		}]
	},
	{
		"cms": "phpshe",
		"rules": [{
			"type": "body",
			"keyword": "content=\"phpshe"
		}]
	},
	{
		"cms": "phpvod",
		"rules": [{
			"type": "body",
			"keyword": "content=\"phpvod"
		}]
	},
	{
		"cms": "phpwind",
		"rules": [{
			"type": "body",
			"keyword": "content=\"phpwind"
		}]
	},
	{
		"cms": "ASPilot-Cart",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Pilot Cart"
		}]
	},
	{
		"cms": "ComersusCart",
		"rules": [{
			"type": "body",
			"keyword": "CONTENT=\"Powered by Comersus"
		}]
	},
	{
		"cms": "DeluxeBB",
		"rules": [{
			"type": "body",
			"keyword": "content=\"powered by DeluxeBB"
		}]
	},
	{
		"cms": "Jcow",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Powered by Jcow"
		}]
	},
	{
		"cms": "PretsaShop",
		"rules": [{
			"type": "body",
			"keyword": "content=\"PrestaShop\""
		}]
	},
	{
		"cms": "richmail",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Richmail"
		}]
	},
	{
		"cms": "FCMS",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Ryan Haudenschilt"
		}]
	},
	{
		"cms": "海洋CMS",
		"rules": [{
			"type": "body",
			"keyword": "content=\"seacms"
		}]
	},
	{
		"cms": "shopbuilder",
		"rules": [{
			"type": "body",
			"keyword": "content=\"ShopBuilder"
		}]
	},
	{
		"cms": "shopex",
		"rules": [{
			"type": "body",
			"keyword": "content=\"ShopEx"
		}]
	},
	{
		"cms": "shopnc",
		"rules": [{
			"type": "body",
			"keyword": "content=\"ShopNC"
		}]
	},
	{
		"cms": "SilverStripe",
		"rules": [{
			"type": "body",
			"keyword": "content=\"SilverStripe"
		}]
	},
	{
		"cms": "Telerik Sitefinity",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Sitefinity"
		}]
	},
	{
		"cms": "STAR CMS",
		"rules": [{
			"type": "body",
			"keyword": "content=\"STARCMS"
		}]
	},
	{
		"cms": "STcms",
		"rules": [{
			"type": "body",
			"keyword": "content=\"STCMS"
		}]
	},
	{
		"cms": "Jamroom",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Talldude Networks"
		}]
	},
	{
		"cms": "FileNice",
		"rules": [{
			"type": "body",
			"keyword": "content=\"the fantabulous mechanical eviltwin code machine"
		}]
	},
	{
		"cms": "Tipask",
		"rules": [{
			"type": "body",
			"keyword": "content=\"tipask"
		}]
	},
	{
		"cms": "泰信TMailer邮件系统",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Tmailer"
		}]
	},
	{
		"cms": "tutucms",
		"rules": [{
			"type": "body",
			"keyword": "content=\"TUTUCMS"
		}]
	},
	{
		"cms": "v5shop",
		"rules": [{
			"type": "body",
			"keyword": "content=\"V5shop"
		}]
	},
	{
		"cms": "Vicworl",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Vicworl"
		}]
	},
	{
		"cms": ".NET",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Visual Basic .NET 7.1"
		}]
	},
	{
		"cms": "AvantFAX",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Web 2.0 HylaFAX"
		}]
	},
	{
		"cms": "CMS-WebManager-Pro",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Webmanager-pro"
		}]
	},
	{
		"cms": "weiphp",
		"rules": [{
			"type": "body",
			"keyword": "content=\"WeiPHP"
		}]
	},
	{
		"cms": "wuzhicms",
		"rules": [{
			"type": "body",
			"keyword": "content=\"wuzhicms"
		}]
	},
	{
		"cms": "BestShopPro",
		"rules": [{
			"type": "body",
			"keyword": "content=\"www.bst.pl"
		}]
	},
	{
		"cms": "CMSQLite",
		"rules": [{
			"type": "body",
			"keyword": "content=\"www.CMSQLite.net"
		}]
	},
	{
		"cms": "易企CMS",
		"rules": [{
			"type": "body",
			"keyword": "content=\"YiqiCMS"
		}]
	},
	{
		"cms": "Yxcms",
		"rules": [{
			"type": "body",
			"keyword": "content=\"Yxcms"
		}]
	},
	{
		"cms": "大米CMS",
		"rules": [{
			"type": "body",
			"keyword": "content=\"大米CMS"
		}]
	},
	{
		"cms": "单点CRM系统",
		"rules": [{
			"type": "body",
			"keyword": "content=\"单点CRM系统"
		}]
	},
	{
		"cms": "dtcms",
		"rules": [{
			"type": "body",
			"keyword": "content=\"动力启航,DTCMS"
		}]
	},
	{
		"cms": "凡科",
		"rules": [{
			"type": "body",
			"keyword": "content=\"凡科"
		}]
	},
	{
		"cms": "汉码软件",
		"rules": [{
			"type": "body",
			"keyword": "content=\"汉码软件"
		}]
	},
	{
		"cms": "74cms",
		"rules": [{
			"type": "body",
			"keyword": "content=\"骑士CMS"
		}]
	},
	{
		"cms": "智睿软件",
		"rules": [{
			"type": "body",
			"keyword": "content=\"智睿软件"
		}]
	},
	{
		"cms": "Intellinet-IP-Camera",
		"rules": [{
			"type": "body",
			"keyword": "Copyright &copy;  INTELLINET NETWORK SOLUTIONS"
		}]
	},
	{
		"cms": "SNB股票交易软件",
		"rules": [{
			"type": "body",
			"keyword": "Copyright 2005–2009 <a href=\"http://www.s-mo.com\">"
		}]
	},
	{
		"cms": "shopnc",
		"rules": [{
			"type": "body",
			"keyword": "Copyright 2007-2014 ShopNC Inc"
		}]
	},
	{
		"cms": "Artiphp-CMS",
		"rules": [{
			"type": "body",
			"keyword": "copyright Artiphp"
		}]
	},
	{
		"cms": "iGENUS邮件系统",
		"rules": [{
			"type": "body",
			"keyword": "Copyright by<A HREF=\"http://www.igenus.org"
		}]
	},
	{
		"cms": "FineCMS",
		"rules": [{
			"type": "body",
			"keyword": "Copyright\" content=\"FineCMS"
		}]
	},
	{
		"cms": "GenieATM",
		"rules": [{
			"type": "body",
			"keyword": "Copyright© Genie Networks Ltd."
		}]
	},
	{
		"cms": "CruxCMS",
		"rules": [{
			"type": "body",
			"keyword": "Created by CruxCMS"
		}]
	},
	{
		"cms": "Foosun",
		"rules": [{
			"type": "body",
			"keyword": "Created by DotNetCMS"
		}]
	},
	{
		"cms": "InterRed",
		"rules": [{
			"type": "body",
			"keyword": "Created with InterRed"
		}]
	},
	{
		"cms": "任我行CRM",
		"rules": [{
			"type": "body",
			"keyword": "CRM_LASTLOGINUSERKEY"
		}]
	},
	{
		"cms": "Django",
		"rules": [{
			"type": "body",
			"keyword": "csrfmiddlewaretoken"
		}]
	},
	{
		"cms": "25yi",
		"rules": [{
			"type": "body",
			"keyword": "css/25yi.css"
		}]
	},
	{
		"cms": "万户网络",
		"rules": [{
			"type": "body",
			"keyword": "css/css_whir.css"
		}]
	},
	{
		"cms": "FestOS",
		"rules": [{
			"type": "body",
			"keyword": "css/festos.css"
		}]
	},
	{
		"cms": "i@Report",
		"rules": [{
			"type": "body",
			"keyword": "css/ireport.css"
		}]
	},
	{
		"cms": "网神防火墙",
		"rules": [{
			"type": "body",
			"keyword": "css/lsec/login.css"
		}]
	},
	{
		"cms": "STcms",
		"rules": [{
			"type": "body",
			"keyword": "DahongY<dahongy@gmail.com>"
		}]
	},
	{
		"cms": "PhpCMS",
		"rules": [{
			"type": "body",
			"keyword": "data/config.js"
		}]
	},
	{
		"cms": "FineCMS",
		"rules": [{
			"type": "body",
			"keyword": "dayrui@gmail.com"
		}]
	},
	{
		"cms": "D-Link-Network-Camera",
		"rules": [{
			"type": "body",
			"keyword": "DCS-950G\".toLowerCase()"
		}]
	},
	{
		"cms": "asp168欧虎",
		"rules": [{
			"type": "body",
			"keyword": "default.php?mod=article&do=detail&tid"
		}]
	},
	{
		"cms": "Apabi数字资源平台",
		"rules": [{
			"type": "body",
			"keyword": "Default/apabi.css"
		}]
	},
	{
		"cms": "GenieATM",
		"rules": [{
			"type": "body",
			"keyword": "defect 3531"
		}]
	},
	{
		"cms": "Daffodil-CRM",
		"rules": [{
			"type": "body",
			"keyword": "Design & Development by Daffodil Software Ltd"
		}]
	},
	{
		"cms": "SOMOIDEA",
		"rules": [{
			"type": "body",
			"keyword": "DESIGN BY SOMOIDEA"
		}]
	},
	{
		"cms": "destoon",
		"rules": [{
			"type": "body",
			"keyword": "destoon_moduleid"
		}]
	},
	{
		"cms": "BrewBlogger",
		"rules": [{
			"type": "body",
			"keyword": "developed by <a href=\"http://www.zkdigital.com"
		}]
	},
	{
		"cms": "易点CMS",
		"rules": [{
			"type": "body",
			"keyword": "DianCMS_SiteName"
		}]
	},
	{
		"cms": "易点CMS",
		"rules": [{
			"type": "body",
			"keyword": "DianCMS_用户登陆引用"
		}]
	},
	{
		"cms": "DnP Firewall",
		"rules": [{
			"type": "body",
			"keyword": "dnp_firewall_redirect"
		}]
	},
	{
		"cms": "校园卡管理系统",
		"rules": [{
			"type": "body",
			"keyword": "document.FormPostds.action=\"xxsearch.action"
		}]
	},
	{
		"cms": "Dolibarr",
		"rules": [{
			"type": "body",
			"keyword": "Dolibarr Development Team"
		}]
	},
	{
		"cms": "Dolibarr",
		"rules": [{
			"type": "body",
			"keyword": "Dolibarr Development Team"
		}]
	},
	{
		"cms": "Dolibarr",
		"rules": [{
			"type": "body",
			"keyword": "Dolibarr Development Team"
		}]
	},
	{
		"cms": "DUgallery",
		"rules": [{
			"type": "body",
			"keyword": "DUgallery"
		}]
	},
	{
		"cms": "DuomiCms",
		"rules": [{
			"type": "body",
			"keyword": "DuomiCms"
		}]
	},
	{
		"cms": "DVWA",
		"rules": [{
			"type": "body",
			"keyword": "dvwa/css/login.css"
		}]
	},
	{
		"cms": "DVWA",
		"rules": [{
			"type": "body",
			"keyword": "dvwa/images/login_logo.png"
		}]
	},
	{
		"cms": "AirvaeCommerce",
		"rules": [{
			"type": "body",
			"keyword": "E-Commerce Shopping Cart Software"
		}]
	},
	{
		"cms": "E-Manage-MySchool",
		"rules": [{
			"type": "body",
			"keyword": "E-Manage All Rights Reserved MySchool Version"
		}]
	},
	{
		"cms": "金蝶EAS",
		"rules": [{
			"type": "body",
			"keyword": "easSessionId"
		}]
	},
	{
		"cms": "ecwapoa",
		"rules": [{
			"type": "body",
			"keyword": "ecwapoa"
		}]
	},
	{
		"cms": "soeasy网站集群系统",
		"rules": [{
			"type": "body",
			"keyword": "EGSS_User"
		}]
	},
	{
		"cms": "eMeeting-Online-Dating-Software",
		"rules": [{
			"type": "body",
			"keyword": "eMeeting Dating Software"
		}]
	},
	{
		"cms": "Jcow",
		"rules": [{
			"type": "body",
			"keyword": "end jcow_application_box"
		}]
	},
	{
		"cms": "EnterCRM",
		"rules": [{
			"type": "body",
			"keyword": "EnterCRM"
		}]
	},
	{
		"cms": "AJA-Video-Converter",
		"rules": [{
			"type": "body",
			"keyword": "eParamID_SWVersion"
		}]
	},
	{
		"cms": "Epiware",
		"rules": [{
			"type": "body",
			"keyword": "Epiware - Project and Document Management"
		}]
	},
	{
		"cms": "i@Report",
		"rules": [{
			"type": "body",
			"keyword": "ESENSOFT_IREPORT_SERVER"
		}]
	},
	{
		"cms": "eSitesBuilder",
		"rules": [{
			"type": "body",
			"keyword": "eSitesBuilder. All rights reserved"
		}]
	},
	{
		"cms": "Etano",
		"rules": [{
			"type": "body",
			"keyword": "Etano</a>. All Rights Reserved."
		}]
	},
	{
		"cms": "SCADA PLC",
		"rules": [{
			"type": "body",
			"keyword": "Ethernet Processor"
		}]
	},
	{
		"cms": "EZCMS",
		"rules": [{
			"type": "body",
			"keyword": "EZCMS Content Management System"
		}]
	},
	{
		"cms": "ezOFFICE",
		"rules": [{
			"type": "body",
			"keyword": "EZOFFICEUSERNAME"
		}]
	},
	{
		"cms": "Fastly cdn",
		"rules": [{
			"type": "body",
			"keyword": "fastcdn.org"
		}]
	},
	{
		"cms": "易瑞授权访问系统",
		"rules": [{
			"type": "body",
			"keyword": "FE0174BB-F093-42AF-AB20-7EC621D10488"
		}]
	},
	{
		"cms": "AVCON6",
		"rules": [{
			"type": "body",
			"keyword": "filename=AVCON6Setup.exe"
		}]
	},
	{
		"cms": "FileNice",
		"rules": [{
			"type": "body",
			"keyword": "fileNice/fileNice.js"
		}]
	},
	{
		"cms": "Foosun",
		"rules": [{
			"type": "body",
			"keyword": "For Foosun"
		}]
	},
	{
		"cms": "ZyXEL",
		"rules": [{
			"type": "body",
			"keyword": "Forms/rpAuth_1"
		}]
	},
	{
		"cms": "FortiGuard",
		"rules": [{
			"type": "body",
			"keyword": "FortiGuard Web Filtering"
		}]
	},
	{
		"cms": "FoxPHP",
		"rules": [{
			"type": "body",
			"keyword": "FoxPHP_ImList"
		}]
	},
	{
		"cms": "FoxPHP",
		"rules": [{
			"type": "body",
			"keyword": "FoxPHPScroll"
		}]
	},
	{
		"cms": "EarlyImpact-ProductCart",
		"rules": [{
			"type": "body",
			"keyword": "fpassword.asp?redirectUrl=&frURL=Custva.asp"
		}]
	},
	{
		"cms": "锐捷NBR路由器",
		"rules": [{
			"type": "body",
			"keyword": "free_nbr_login_form.png"
		}]
	},
	{
		"cms": "锐捷NBR路由器",
		"rules": [{
			"type": "body",
			"keyword": "free_nbr_login_form.png"
		}]
	},
	{
		"cms": "e-junkie",
		"rules": [{
			"type": "body",
			"keyword": "function EJEJC_lc"
		}]
	},
	{
		"cms": "esoTalk",
		"rules": [{
			"type": "body",
			"keyword": "generated by esoTalk"
		}]
	},
	{
		"cms": "phpDocumentor",
		"rules": [{
			"type": "body",
			"keyword": "Generated by phpDocumentor"
		}]
	},
	{
		"cms": "phpDocumentor",
		"rules": [{
			"type": "body",
			"keyword": "Generated by phpDocumentor"
		}]
	},
	{
		"cms": "1024cms",
		"rules": [{
			"type": "body",
			"keyword": "generator\" content=\"1024 CMS (c)"
		}]
	},
	{
		"cms": "2z project",
		"rules": [{
			"type": "body",
			"keyword": "Generator\" content=\"2z project"
		}]
	},
	{
		"cms": "6kbbs",
		"rules": [{
			"type": "body",
			"keyword": "generator\" content=\"6KBBS"
		}]
	},
	{
		"cms": "6kbbs",
		"rules": [{
			"type": "body",
			"keyword": "generator\" content=\"6KBBS"
		}]
	},
	{
		"cms": "Adobe_GoLive",
		"rules": [{
			"type": "body",
			"keyword": "generator\" content=\"Adobe GoLive"
		}]
	},
	{
		"cms": "Adobe_RoboHelp",
		"rules": [{
			"type": "body",
			"keyword": "generator\" content=\"Adobe RoboHelp"
		}]
	},
	{
		"cms": "Amaya",
		"rules": [{
			"type": "body",
			"keyword": "generator\" content=\"Amaya"
		}]
	},
	{
		"cms": "awstats_admin",
		"rules": [{
			"type": "body",
			"keyword": "generator\" content=\"AWStats"
		}]
	},
	{
		"cms": "Centreon",
		"rules": [{
			"type": "body",
			"keyword": "Generator\" content=\"Centreon - Copyright"
		}]
	},
	{
		"cms": "MediaWiki",
		"rules": [{
			"type": "body",
			"keyword": "generator\" content=\"MediaWiki"
		}]
	},
	{
		"cms": "Typecho",
		"rules": [{
			"type": "body",
			"keyword": "generator\" content=\"Typecho"
		}]
	},
	{
		"cms": "webEdition",
		"rules": [{
			"type": "body",
			"keyword": "generator\" content=\"webEdition"
		}]
	},
	{
		"cms": "用友U8",
		"rules": [{
			"type": "body",
			"keyword": "getFirstU8Accid"
		}]
	},
	{
		"cms": "GridSite",
		"rules": [{
			"type": "body",
			"keyword": "gridsite-admin.cgi?cmd"
		}]
	},
	{
		"cms": " 海盗云商(Haidao)",
		"rules": [{
			"type": "body",
			"keyword": "haidao.web.general.js"
		}]
	},
	{
		"cms": "校园卡管理系统",
		"rules": [{
			"type": "body",
			"keyword": "Harbin synjones electronic"
		}]
	},
	{
		"cms": "Kayako-SupportSuite",
		"rules": [{
			"type": "body",
			"keyword": "Help Desk Software By Kayako eSupport"
		}]
	},
	{
		"cms": "HESK",
		"rules": [{
			"type": "body",
			"keyword": "hesk_javascript.js"
		}]
	},
	{
		"cms": "HESK",
		"rules": [{
			"type": "body",
			"keyword": "hesk_style.css"
		}]
	},
	{
		"cms": "HUAWEI Inner Web",
		"rules": [{
			"type": "body",
			"keyword": "hidden_frame.html"
		}]
	},
	{
		"cms": "HIMS酒店云计算服务",
		"rules": [{
			"type": "body",
			"keyword": "HIMS酒店云计算服务"
		}]
	},
	{
		"cms": "hishop",
		"rules": [{
			"type": "body",
			"keyword": "Hishop development team"
		}]
	},
	{
		"cms": "hishop",
		"rules": [{
			"type": "body",
			"keyword": "hishop.plugins.openid"
		}]
	},
	{
		"cms": "擎天电子政务",
		"rules": [{
			"type": "body",
			"keyword": "homepages/content_page.aspx"
		}]
	},
	{
		"cms": "云因网上书店",
		"rules": [{
			"type": "body",
			"keyword": "href=\"../css/newscomm.css"
		}]
	},
	{
		"cms": "GeoNode",
		"rules": [{
			"type": "body",
			"keyword": "href=\"/catalogue/opensearch\" title=\"GeoNode Search"
		}]
	},
	{
		"cms": "泰信TMailer邮件系统",
		"rules": [{
			"type": "body",
			"keyword": "href=\"/tmailer/img/logo/favicon.ico"
		}]
	},
	{
		"cms": "锐商企业CMS",
		"rules": [{
			"type": "body",
			"keyword": "href=\"/Writable/ClientImages/mycss.css"
		}]
	},
	{
		"cms": "BugTracker.NET",
		"rules": [{
			"type": "body",
			"keyword": "href=\"btnet.css"
		}]
	},
	{
		"cms": "ComersusCart",
		"rules": [{
			"type": "body",
			"keyword": "href=\"comersus_showCart.asp"
		}]
	},
	{
		"cms": "Gossamer-Forum",
		"rules": [{
			"type": "body",
			"keyword": "href=\"gforum.cgi?username="
		}]
	},
	{
		"cms": "Cachelogic-Expired-Domains-Script",
		"rules": [{
			"type": "body",
			"keyword": "href=\"http://cachelogic.net\">Cachelogic.net"
		}]
	},
	{
		"cms": "ClipBucket",
		"rules": [{
			"type": "body",
			"keyword": "href=\"http://clip-bucket.com/\">ClipBucket"
		}]
	},
	{
		"cms": "CMS-WebManager-Pro",
		"rules": [{
			"type": "body",
			"keyword": "href=\"http://webmanager-pro.com\">Web.Manager"
		}]
	},
	{
		"cms": "bitweaver",
		"rules": [{
			"type": "body",
			"keyword": "href=\"http://www.bitweaver.org\">Powered by"
		}]
	},
	{
		"cms": "Imageview",
		"rules": [{
			"type": "body",
			"keyword": "href=\"http://www.blackdot.be\" title=\"Blackdot.be"
		}]
	},
	{
		"cms": "iTop",
		"rules": [{
			"type": "body",
			"keyword": "href=\"http://www.combodo.com/itop"
		}]
	},
	{
		"cms": "Dokeos",
		"rules": [{
			"type": "body",
			"keyword": "href=\"http://www.dokeos.com\" rel=\"Copyright"
		}]
	},
	{
		"cms": "iLO",
		"rules": [{
			"type": "body",
			"keyword": "href=\"http://www.hp.com/go/ilo"
		}]
	},
	{
		"cms": "JagoanStore",
		"rules": [{
			"type": "body",
			"keyword": "href=\"http://www.jagoanstore.com/\" target=\"_blank\">Toko Online"
		}]
	},
	{
		"cms": "FrogCMS",
		"rules": [{
			"type": "body",
			"keyword": "href=\"http://www.madebyfrog.com\">Frog CMS"
		}]
	},
	{
		"cms": "FormMail",
		"rules": [{
			"type": "body",
			"keyword": "href=\"http://www.worldwidemart.com/scripts/formmail.shtml"
		}]
	},
	{
		"cms": "i-Gallery",
		"rules": [{
			"type": "body",
			"keyword": "href=\"igallery.asp"
		}]
	},
	{
		"cms": "JGS-Portal",
		"rules": [{
			"type": "body",
			"keyword": "href=\"jgs_portal_box.php?id="
		}]
	},
	{
		"cms": "ipTIME-Router",
		"rules": [{
			"type": "body",
			"keyword": "href=iptime.css"
		}]
	},
	{
		"cms": "Kampyle",
		"rules": [{
			"type": "body",
			"keyword": "http://cf.kampyle.com/k_button.js"
		}]
	},
	{
		"cms": "小脑袋",
		"rules": [{
			"type": "body",
			"keyword": "http://stat.xiaonaodai.com/stat.php"
		}]
	},
	{
		"cms": "Barracuda-Spam-Firewall",
		"rules": [{
			"type": "body",
			"keyword": "http://www.barracudanetworks.com?a=bsf_product"
		}]
	},
	{
		"cms": "Claroline",
		"rules": [{
			"type": "body",
			"keyword": "http://www.claroline.net\" rel=\"Copyright"
		}]
	},
	{
		"cms": "Intellinet-IP-Camera",
		"rules": [{
			"type": "body",
			"keyword": "http://www.intellinet-network.com/driver/NetCam.exe"
		}]
	},
	{
		"cms": "信达OA",
		"rules": [{
			"type": "body",
			"keyword": "http://www.xdoa.cn</a>"
		}]
	},
	{
		"cms": "宝塔面板-python",
		"rules": [{
			"type": "body",
			"keyword": "https://www.bt.cn/bbs/thread-18367-1-1.html"
		}]
	},
	{
		"cms": "海天OA",
		"rules": [{
			"type": "body",
			"keyword": "HTVOS.js"
		}]
	},
	{
		"cms": "科来RAS",
		"rules": [{
			"type": "body",
			"keyword": "i18ninit.min.js"
		}]
	},
	{
		"cms": "H3C公司产品",
		"rules": [{
			"type": "body",
			"keyword": "icg_helpScript.js"
		}]
	},
	{
		"cms": "JXT-Consulting",
		"rules": [{
			"type": "body",
			"keyword": "id=\"jxt-popup-wrapper"
		}]
	},
	{
		"cms": "BugFree",
		"rules": [{
			"type": "body",
			"keyword": "id=\"logo\" alt=BugFree"
		}]
	},
	{
		"cms": "AvantFAX",
		"rules": [{
			"type": "body",
			"keyword": "images/avantfax-big.png"
		}]
	},
	{
		"cms": "winwebmail",
		"rules": [{
			"type": "body",
			"keyword": "images/owin.css"
		}]
	},
	{
		"cms": "VOS3000",
		"rules": [{
			"type": "body",
			"keyword": "images/vos3000.ico"
		}]
	},
	{
		"cms": "xoops",
		"rules": [{
			"type": "body",
			"keyword": "include/xoops.js"
		}]
	},
	{
		"cms": "AppServ",
		"rules": [{
			"type": "body",
			"keyword": "index.php?appservlang=th"
		}]
	},
	{
		"cms": "Auxilium-PetRatePro",
		"rules": [{
			"type": "body",
			"keyword": "index.php?cmd=11"
		}]
	},
	{
		"cms": "ClanSphere",
		"rules": [{
			"type": "body",
			"keyword": "index.php?mod=clansphere&amp;action=about"
		}]
	},
	{
		"cms": "infoglue",
		"rules": [{
			"type": "body",
			"keyword": "infoglueBox.png"
		}]
	},
	{
		"cms": "sony摄像头",
		"rules": [{
			"type": "body",
			"keyword": "inquiry.cgi?inqjs=system&inqjs=camera"
		}]
	},
	{
		"cms": "eagleeyescctv",
		"rules": [{
			"type": "body",
			"keyword": "IP Surveillance for Your Life"
		}]
	},
	{
		"cms": "IP.Board",
		"rules": [{
			"type": "body",
			"keyword": "ipb.vars"
		}]
	},
	{
		"cms": "i@Report",
		"rules": [{
			"type": "body",
			"keyword": "ireportclient"
		}]
	},
	{
		"cms": "bbPress",
		"rules": [{
			"type": "body",
			"keyword": "is proudly powered by <a href=\"http://bbpress.org"
		}]
	},
	{
		"cms": "Lotus",
		"rules": [{
			"type": "body",
			"keyword": "iwaredir.nsf"
		}]
	},
	{
		"cms": "Lotus",
		"rules": [{
			"type": "body",
			"keyword": "iwaredir.nsf"
		}]
	},
	{
		"cms": "javashop",
		"rules": [{
			"type": "body",
			"keyword": "javashop微信公众号"
		}]
	},
	{
		"cms": "jCore",
		"rules": [{
			"type": "body",
			"keyword": "JCORE_VERSION = "
		}]
	},
	{
		"cms": "jobberBase",
		"rules": [{
			"type": "body",
			"keyword": "Jobber.PerformSearch"
		}]
	},
	{
		"cms": "Tiki-wiki CMS",
		"rules": [{
			"type": "body",
			"keyword": "jqueryTiki = new Object"
		}]
	},
	{
		"cms": "HP_iLO(HP_Integrated_Lights-Out)",
		"rules": [{
			"type": "body",
			"keyword": "js/iLO.js"
		}]
	},
	{
		"cms": "H3C-SecBlade-FireWall",
		"rules": [{
			"type": "body",
			"keyword": "js/MulPlatAPI.js"
		}]
	},
	{
		"cms": "Kibana",
		"rules": [{
			"type": "body",
			"keyword": "kbnVersion"
		}]
	},
	{
		"cms": "科蚁CMS",
		"rules": [{
			"type": "body",
			"keyword": "keyicms：keyicms"
		}]
	},
	{
		"cms": "地平线CMS",
		"rules": [{
			"type": "body",
			"keyword": "labelOppInforStyle"
		}]
	},
	{
		"cms": "易普拉格科研管理系统",
		"rules": [{
			"type": "body",
			"keyword": "lan12-jingbian-hong"
		}]
	},
	{
		"cms": "浪潮政务系统",
		"rules": [{
			"type": "body",
			"keyword": "LangChao.ECGAP.OutPortal"
		}]
	},
	{
		"cms": "AVCON6",
		"rules": [{
			"type": "body",
			"keyword": "language_dispose.action"
		}]
	},
	{
		"cms": "lemis管理系统",
		"rules": [{
			"type": "body",
			"keyword": "lemis.WEB_APP_NAME"
		}]
	},
	{
		"cms": "78oa",
		"rules": [{
			"type": "body",
			"keyword": "license.78oa.com"
		}]
	},
	{
		"cms": "科信邮件系统",
		"rules": [{
			"type": "body",
			"keyword": "lo_computername"
		}]
	},
	{
		"cms": "IQeye-Netcam",
		"rules": [{
			"type": "body",
			"keyword": "loc = \"iqeyevid.html"
		}]
	},
	{
		"cms": "金龙卡金融化一卡通网站查询子系统",
		"rules": [{
			"type": "body",
			"keyword": "location.href=\"homeLogin.action"
		}]
	},
	{
		"cms": "OpenMas",
		"rules": [{
			"type": "body",
			"keyword": "loginHead\"><link href=\"App_Themes"
		}]
	},
	{
		"cms": "OpenMas",
		"rules": [{
			"type": "body",
			"keyword": "loginHead\"><link href=\"App_Themes"
		}]
	},
	{
		"cms": "FreeboxOS",
		"rules": [{
			"type": "body",
			"keyword": "logo_freeboxos"
		}]
	},
	{
		"cms": "UFIDA_NC",
		"rules": [{
			"type": "body",
			"keyword": "logo/images/ufida_nc.png"
		}]
	},
	{
		"cms": "UFIDA_NC",
		"rules": [{
			"type": "body",
			"keyword": "logo/images/ufida_nc.png"
		}]
	},
	{
		"cms": "UFIDA_NC",
		"rules": [{
			"type": "body",
			"keyword": "logo/images/ufida_nc.png"
		}]
	},
	{
		"cms": "lynxspring_JENEsys",
		"rules": [{
			"type": "body",
			"keyword": "LX JENEsys"
		}]
	},
	{
		"cms": "IdeaCMS",
		"rules": [{
			"type": "body",
			"keyword": "m_ctr32"
		}]
	},
	{
		"cms": "苹果CMS",
		"rules": [{
			"type": "body",
			"keyword": "maccms:voddaycount"
		}]
	},
	{
		"cms": "Magento",
		"rules": [{
			"type": "body",
			"keyword": "Magento, Varien, E-commerce"
		}]
	},
	{
		"cms": "BASE",
		"rules": [{
			"type": "body",
			"keyword": "mailto:base@secureideas.net"
		}]
	},
	{
		"cms": "云因网上书店",
		"rules": [{
			"type": "body",
			"keyword": "main/building.cfm"
		}]
	},
	{
		"cms": "JBoss_AS",
		"rules": [{
			"type": "body",
			"keyword": "Manage this JBoss AS Instance"
		}]
	},
	{
		"cms": "mantis",
		"rules": [{
			"type": "body",
			"keyword": "MantisBT Team"
		}]
	},
	{
		"cms": "Maticsoft_Shop_动软商城",
		"rules": [{
			"type": "body",
			"keyword": "Maticsoft Shop"
		}]
	},
	{
		"cms": "MaticsoftSNS_动软分享社区",
		"rules": [{
			"type": "body",
			"keyword": "MaticsoftSNS"
		}]
	},
	{
		"cms": "华为 MCU",
		"rules": [{
			"type": "body",
			"keyword": "McuR5-min.js"
		}]
	},
	{
		"cms": "华为 MCU",
		"rules": [{
			"type": "body",
			"keyword": "MCUType.js"
		}]
	},
	{
		"cms": "网动云视讯平台",
		"rules": [{
			"type": "body",
			"keyword": "meetingShow!show.action"
		}]
	},
	{
		"cms": "moosefs",
		"rules": [{
			"type": "body",
			"keyword": "mfs.cgi"
		}]
	},
	{
		"cms": "moosefs",
		"rules": [{
			"type": "body",
			"keyword": "mfs.cgi"
		}]
	},
	{
		"cms": "tp-shop",
		"rules": [{
			"type": "body",
			"keyword": "mn-c-top"
		}]
	},
	{
		"cms": "帕拉迪统一安全管理和综合审计系统",
		"rules": [{
			"type": "body",
			"keyword": "module/image/pldsec.css"
		}]
	},
	{
		"cms": "Ruckus",
		"rules": [{
			"type": "body",
			"keyword": "mon.  Tell me your username"
		}]
	},
	{
		"cms": "Panasonic Network Camera",
		"rules": [{
			"type": "body",
			"keyword": "MultiCameraFrame?Mode=Motion&Language"
		}]
	},
	{
		"cms": "Munin",
		"rules": [{
			"type": "body",
			"keyword": "munin-month.html"
		}]
	},
	{
		"cms": "BugFree",
		"rules": [{
			"type": "body",
			"keyword": "name=\"BugUserPWD"
		}]
	},
	{
		"cms": "DublinCore",
		"rules": [{
			"type": "body",
			"keyword": "name=\"DC.title"
		}]
	},
	{
		"cms": "cApexWEB",
		"rules": [{
			"type": "body",
			"keyword": "name=\"dfparentdb"
		}]
	},
	{
		"cms": "DnP-Firewall",
		"rules": [{
			"type": "body",
			"keyword": "name=\"dnp_firewall_redirect"
		}]
	},
	{
		"cms": "Apache-Forrest",
		"rules": [{
			"type": "body",
			"keyword": "name=\"Forrest"
		}]
	},
	{
		"cms": "Dokeos",
		"rules": [{
			"type": "body",
			"keyword": "name=\"Generator\" content=\"Dokeos"
		}]
	},
	{
		"cms": "fckeditor",
		"rules": [{
			"type": "body",
			"keyword": "new FCKeditor"
		}]
	},
	{
		"cms": "NITC",
		"rules": [{
			"type": "body",
			"keyword": "NITC Web Marketing Service"
		}]
	},
	{
		"cms": "NetDvrV3",
		"rules": [{
			"type": "body",
			"keyword": "objLvrForNoIE"
		}]
	},
	{
		"cms": "通达OA",
		"rules": [{
			"type": "body",
			"keyword": "Office Anywhere 2013"
		}]
	},
	{
		"cms": "天融信网络审计系统",
		"rules": [{
			"type": "body",
			"keyword": "onclick=\"dlg_download()"
		}]
	},
	{
		"cms": "浪潮政务系统",
		"rules": [{
			"type": "body",
			"keyword": "OnlineQuery/QueryList.aspx"
		}]
	},
	{
		"cms": "北创图书检索系统",
		"rules": [{
			"type": "body",
			"keyword": "opac_two"
		}]
	},
	{
		"cms": "北创图书检索系统",
		"rules": [{
			"type": "body",
			"keyword": "opac_two"
		}]
	},
	{
		"cms": "Oracle_OPERA",
		"rules": [{
			"type": "body",
			"keyword": "OperaLogin/Welcome.do"
		}]
	},
	{
		"cms": "oracle_applicaton_server",
		"rules": [{
			"type": "body",
			"keyword": "OraLightHeaderSub"
		}]
	},
	{
		"cms": "Parallels Plesk Panel",
		"rules": [{
			"type": "body",
			"keyword": "Parallels IP Holdings GmbH"
		}]
	},
	{
		"cms": "arrisi_Touchstone",
		"rules": [{
			"type": "body",
			"keyword": "passWithWarnings"
		}]
	},
	{
		"cms": "phpweb",
		"rules": [{
			"type": "body",
			"keyword": "PDV_PAGENAME"
		}]
	},
	{
		"cms": "PHPMyWind",
		"rules": [{
			"type": "body",
			"keyword": "phpMyWind.com All Rights Reserved"
		}]
	},
	{
		"cms": "易分析",
		"rules": [{
			"type": "body",
			"keyword": "PHPStat Analytics 网站数据分析系统"
		}]
	},
	{
		"cms": "BlogEngine_NET",
		"rules": [{
			"type": "body",
			"keyword": "pics/blogengine.ico"
		}]
	},
	{
		"cms": "milu_seotool",
		"rules": [{
			"type": "body",
			"keyword": "plugin.php?id=milu_seotool"
		}]
	},
	{
		"cms": "bluecms",
		"rules": [{
			"type": "body",
			"keyword": "power by bcms"
		}]
	},
	{
		"cms": "CuuMall",
		"rules": [{
			"type": "body",
			"keyword": "Power by CuuMall"
		}]
	},
	{
		"cms": "DedeCMS",
		"rules": [{
			"type": "body",
			"keyword": "Power by DedeCms"
		}]
	},
	{
		"cms": "doccms",
		"rules": [{
			"type": "body",
			"keyword": "Power by DocCms"
		}]
	},
	{
		"cms": "appcms",
		"rules": [{
			"type": "body",
			"keyword": "Powerd by AppCMS"
		}]
	},
	{
		"cms": "BoonEx-Dolphin",
		"rules": [{
			"type": "body",
			"keyword": "Powered by                    Dolphin - <a href=\"http://www.boonex.com/products/dolphin"
		}]
	},
	{
		"cms": "Energine",
		"rules": [{
			"type": "body",
			"keyword": "Powered by <a href=\"http://energine.org/"
		}]
	},
	{
		"cms": "boastMachine",
		"rules": [{
			"type": "body",
			"keyword": "Powered by <a href=\"http://boastology.com"
		}]
	},
	{
		"cms": "CF-Image-Hosting-Script",
		"rules": [{
			"type": "body",
			"keyword": "Powered By <a href=\"http://codefuture.co.uk/projects/imagehost/"
		}]
	},
	{
		"cms": "F3Site",
		"rules": [{
			"type": "body",
			"keyword": "Powered by <a href=\"http://compmaster.prv.pl"
		}]
	},
	{
		"cms": "Dotclear",
		"rules": [{
			"type": "body",
			"keyword": "Powered by <a href=\"http://dotclear.org/"
		}]
	},
	{
		"cms": "FluxBB",
		"rules": [{
			"type": "body",
			"keyword": "Powered by <a href=\"http://fluxbb.org/"
		}]
	},
	{
		"cms": "chillyCMS",
		"rules": [{
			"type": "body",
			"keyword": "powered by <a href=\"http://FrozenPepper.de"
		}]
	},
	{
		"cms": "GeoNode",
		"rules": [{
			"type": "body",
			"keyword": "Powered by <a href=\"http://geonode.org"
		}]
	},
	{
		"cms": "HostBill",
		"rules": [{
			"type": "body",
			"keyword": "Powered by <a href=\"http://hostbillapp.com"
		}]
	},
	{
		"cms": "iScripts-MultiCart",
		"rules": [{
			"type": "body",
			"keyword": "Powered by <a href=\"http://iscripts.com/multicart"
		}]
	},
	{
		"cms": "74cms",
		"rules": [{
			"type": "body",
			"keyword": "Powered by <a href=\"http://www.74cms.com/\""
		}]
	},
	{
		"cms": "AV-Arcade",
		"rules": [{
			"type": "body",
			"keyword": "Powered by <a href=\"http://www.avscripts.net/avarcade/"
		}]
	},
	{
		"cms": "BloofoxCMS",
		"rules": [{
			"type": "body",
			"keyword": "Powered by <a href=\"http://www.bloofox.com"
		}]
	},
	{
		"cms": "CalendarScript",
		"rules": [{
			"type": "body",
			"keyword": "Powered by <A HREF=\"http://www.CalendarScript.com"
		}]
	},
	{
		"cms": "ClipShare",
		"rules": [{
			"type": "body",
			"keyword": "Powered By <a href=\"http://www.clip-share.com"
		}]
	},
	{
		"cms": "Etano",
		"rules": [{
			"type": "body",
			"keyword": "Powered by <a href=\"http://www.datemill.com"
		}]
	},
	{
		"cms": "EasyConsole-CMS",
		"rules": [{
			"type": "body",
			"keyword": "Powered by <a href=\"http://www.easyconsole.com"
		}]
	},
	{
		"cms": "eDirectory",
		"rules": [{
			"type": "body",
			"keyword": "Powered by <a href=\"http://www.edirectory.com"
		}]
	},
	{
		"cms": "EduSoho开源网络课堂",
		"rules": [{
			"type": "body",
			"keyword": "Powered by <a href=\"http://www.edusoho.com"
		}]
	},
	{
		"cms": "cInvoice",
		"rules": [{
			"type": "body",
			"keyword": "Powered by <a href=\"http://www.forperfect.com/"
		}]
	},
	{
		"cms": "GeekLog",
		"rules": [{
			"type": "body",
			"keyword": "Powered By <a href=\"http://www.geeklog.net/"
		}]
	},
	{
		"cms": "HESK",
		"rules": [{
			"type": "body",
			"keyword": "Powered by <a href=\"http://www.hesk.com"
		}]
	},
	{
		"cms": "Hycus-CMS",
		"rules": [{
			"type": "body",
			"keyword": "Powered By <a href=\"http://www.hycus.com"
		}]
	},
	{
		"cms": "Ikonboard",
		"rules": [{
			"type": "body",
			"keyword": "Powered by <a href=\"http://www.ikonboard.com"
		}]
	},
	{
		"cms": "InvisionPowerBoard",
		"rules": [{
			"type": "body",
			"keyword": "Powered by <a href=\"http://www.invisionboard.com"
		}]
	},
	{
		"cms": "ischoolsite",
		"rules": [{
			"type": "body",
			"keyword": "Powered by <a href=\"http://www.ischoolsite.com"
		}]
	},
	{
		"cms": "iScripts-ReserveLogic",
		"rules": [{
			"type": "body",
			"keyword": "Powered by <a href=\"http://www.iscripts.com/reservelogic/"
		}]
	},
	{
		"cms": "ISPConfig",
		"rules": [{
			"type": "body",
			"keyword": "powered by <a HREF=\"http://www.ispconfig.org"
		}]
	},
	{
		"cms": "科蚁CMS",
		"rules": [{
			"type": "body",
			"keyword": "Powered by <a href=\"http://www.keyicms.com"
		}]
	},
	{
		"cms": "Car-Portal",
		"rules": [{
			"type": "body",
			"keyword": "Powered by <a href=\"http://www.netartmedia.net/carsportal"
		}]
	},
	{
		"cms": "Buddy-Zone",
		"rules": [{
			"type": "body",
			"keyword": "Powered By <a href=\"http://www.vastal.com"
		}]
	},
	{
		"cms": "HESK",
		"rules": [{
			"type": "body",
			"keyword": "Powered by <a href=\"https://www.hesk.com"
		}]
	},
	{
		"cms": "Burning-Board-Lite",
		"rules": [{
			"type": "body",
			"keyword": "Powered by <b><a href=\"http://www.woltlab.de"
		}]
	},
	{
		"cms": "Burning-Board-Lite",
		"rules": [{
			"type": "body",
			"keyword": "Powered by <b>Burning Board"
		}]
	},
	{
		"cms": "JGS-Portal",
		"rules": [{
			"type": "body",
			"keyword": "Powered by <b>JGS-Portal Version"
		}]
	},
	{
		"cms": "discuz",
		"rules": [{
			"type": "body",
			"keyword": "Powered by <strong><a href=\"http://www.discuz.net"
		}]
	},
	{
		"cms": "1024cms",
		"rules": [{
			"type": "body",
			"keyword": "Powered by 1024 CMS"
		}]
	},
	{
		"cms": "1024 CMS",
		"rules": [{
			"type": "body",
			"keyword": "Powered by 1024 CMS"
		}]
	},
	{
		"cms": "25yi",
		"rules": [{
			"type": "body",
			"keyword": "Powered by 25yi"
		}]
	},
	{
		"cms": "6kbbs",
		"rules": [{
			"type": "body",
			"keyword": "Powered by 6kbbs"
		}]
	},
	{
		"cms": "6kbbs",
		"rules": [{
			"type": "body",
			"keyword": "Powered by 6kbbs"
		}]
	},
	{
		"cms": "Acidcat_CMS",
		"rules": [{
			"type": "body",
			"keyword": "Powered by Acidcat CMS"
		}]
	},
	{
		"cms": "Acidcat CMS",
		"rules": [{
			"type": "body",
			"keyword": "Powered by Acidcat CMS"
		}]
	},
	{
		"cms": "activeCollab",
		"rules": [{
			"type": "body",
			"keyword": "powered by activeCollab"
		}]
	},
	{
		"cms": "AM4SS",
		"rules": [{
			"type": "body",
			"keyword": "Powered by am4ss"
		}]
	},
	{
		"cms": "Atmail-WebMail",
		"rules": [{
			"type": "body",
			"keyword": "Powered by Atmail"
		}]
	},
	{
		"cms": "Auto-CMS",
		"rules": [{
			"type": "body",
			"keyword": "Powered by Auto CMS"
		}]
	},
	{
		"cms": "b2evolution",
		"rules": [{
			"type": "body",
			"keyword": "Powered by b2evolution"
		}]
	},
	{
		"cms": "boastMachine",
		"rules": [{
			"type": "body",
			"keyword": "powered by boastMachine"
		}]
	},
	{
		"cms": "BrowserCMS",
		"rules": [{
			"type": "body",
			"keyword": "Powered by BrowserCMS"
		}]
	},
	{
		"cms": "Bulletlink-Newspaper-Template",
		"rules": [{
			"type": "body",
			"keyword": "powered by bulletlink"
		}]
	},
	{
		"cms": "CaupoShop-Classic",
		"rules": [{
			"type": "body",
			"keyword": "Powered by CaupoShop"
		}]
	},
	{
		"cms": "CitusCMS",
		"rules": [{
			"type": "body",
			"keyword": "Powered by CitusCMS"
		}]
	},
	{
		"cms": "CMSimple",
		"rules": [{
			"type": "body",
			"keyword": "Powered by CMSimple.dk"
		}]
	},
	{
		"cms": "CMSQLite",
		"rules": [{
			"type": "body",
			"keyword": "powered by CMSQLite"
		}]
	},
	{
		"cms": "协众OA",
		"rules": [{
			"type": "body",
			"keyword": "Powered by CNOA.CN"
		}]
	},
	{
		"cms": "Contrexx-CMS",
		"rules": [{
			"type": "body",
			"keyword": "powered by Contrexx"
		}]
	},
	{
		"cms": "Cyn_in",
		"rules": [{
			"type": "body",
			"keyword": "Powered by cyn.in"
		}]
	},
	{
		"cms": "Daffodil-CRM",
		"rules": [{
			"type": "body",
			"keyword": "Powered by Daffodil"
		}]
	},
	{
		"cms": "DBHcms",
		"rules": [{
			"type": "body",
			"keyword": "powered by DBHcms"
		}]
	},
	{
		"cms": "Diferior",
		"rules": [{
			"type": "body",
			"keyword": "Powered by Diferior"
		}]
	},
	{
		"cms": "DnP Firewall",
		"rules": [{
			"type": "body",
			"keyword": "Powered by DnP Firewall"
		}]
	},
	{
		"cms": "DokuWiki",
		"rules": [{
			"type": "body",
			"keyword": "powered by DokuWiki"
		}]
	},
	{
		"cms": "DouPHP",
		"rules": [{
			"type": "body",
			"keyword": "Powered by DouPHP"
		}]
	},
	{
		"cms": "DrugPak",
		"rules": [{
			"type": "body",
			"keyword": "Powered by DrugPak"
		}]
	},
	{
		"cms": "dswjcms",
		"rules": [{
			"type": "body",
			"keyword": "Powered by Dswjcms"
		}]
	},
	{
		"cms": "DT-Centrepiece",
		"rules": [{
			"type": "body",
			"keyword": "Powered By DT Centrepiece"
		}]
	},
	{
		"cms": "DUgallery",
		"rules": [{
			"type": "body",
			"keyword": "Powered by DUportal"
		}]
	},
	{
		"cms": "EasyConsole-CMS",
		"rules": [{
			"type": "body",
			"keyword": "Powered by EasyConsole CMS"
		}]
	},
	{
		"cms": "E-Xoopport",
		"rules": [{
			"type": "body",
			"keyword": "Powered by E-Xoopport"
		}]
	},
	{
		"cms": "EazyCMS",
		"rules": [{
			"type": "body",
			"keyword": "powered by eazyCMS"
		}]
	},
	{
		"cms": "Echo",
		"rules": [{
			"type": "body",
			"keyword": "powered by echo"
		}]
	},
	{
		"cms": "EduSoho开源网络课堂",
		"rules": [{
			"type": "body",
			"keyword": "Powered By EduSoho"
		}]
	},
	{
		"cms": "Elite-Gaming-Ladders",
		"rules": [{
			"type": "body",
			"keyword": "Powered by Elite"
		}]
	},
	{
		"cms": "AlstraSoft-EPay-Enterprise",
		"rules": [{
			"type": "body",
			"keyword": "Powered by EPay Enterprise"
		}]
	},
	{
		"cms": "esoTalk",
		"rules": [{
			"type": "body",
			"keyword": "Powered by esoTalk"
		}]
	},
	{
		"cms": "ESPCMS",
		"rules": [{
			"type": "body",
			"keyword": "Powered by ESPCMS"
		}]
	},
	{
		"cms": "Esvon-Classifieds",
		"rules": [{
			"type": "body",
			"keyword": "Powered by Esvon"
		}]
	},
	{
		"cms": "eTicket",
		"rules": [{
			"type": "body",
			"keyword": "Powered by eTicket"
		}]
	},
	{
		"cms": "Exponent-CMS",
		"rules": [{
			"type": "body",
			"keyword": "Powered by Exponent CMS"
		}]
	},
	{
		"cms": "EZCMS",
		"rules": [{
			"type": "body",
			"keyword": "Powered by EZCMS"
		}]
	},
	{
		"cms": "FCMS",
		"rules": [{
			"type": "body",
			"keyword": "Powered by Family Connections"
		}]
	},
	{
		"cms": "fengcms",
		"rules": [{
			"type": "body",
			"keyword": "Powered by FengCms"
		}]
	},
	{
		"cms": "FineCMS",
		"rules": [{
			"type": "body",
			"keyword": "Powered by FineCMS"
		}]
	},
	{
		"cms": "Flyspray",
		"rules": [{
			"type": "body",
			"keyword": "Powered by Flyspray"
		}]
	},
	{
		"cms": "GetSimple",
		"rules": [{
			"type": "body",
			"keyword": "Powered by GetSimple"
		}]
	},
	{
		"cms": "HoloCMS",
		"rules": [{
			"type": "body",
			"keyword": "Powered by HoloCMS"
		}]
	},
	{
		"cms": "ICEshop",
		"rules": [{
			"type": "body",
			"keyword": "Powered by ICEshop"
		}]
	},
	{
		"cms": "IdeaCMS",
		"rules": [{
			"type": "body",
			"keyword": "Powered By IdeaCMS"
		}]
	},
	{
		"cms": "IMGCms",
		"rules": [{
			"type": "body",
			"keyword": "Powered by IMGCMS"
		}]
	},
	{
		"cms": "Inout-Adserver",
		"rules": [{
			"type": "body",
			"keyword": "Powered by Inoutscripts"
		}]
	},
	{
		"cms": "JXT-Consulting",
		"rules": [{
			"type": "body",
			"keyword": "Powered by JXT Consulting"
		}]
	},
	{
		"cms": "KaiBB",
		"rules": [{
			"type": "body",
			"keyword": "Powered by KaiBB"
		}]
	},
	{
		"cms": "Kajona",
		"rules": [{
			"type": "body",
			"keyword": "powered by Kajona"
		}]
	},
	{
		"cms": "Kayako-SupportSuite",
		"rules": [{
			"type": "body",
			"keyword": "Powered By Kayako eSupport"
		}]
	},
	{
		"cms": "kingcms",
		"rules": [{
			"type": "body",
			"keyword": "Powered by KingCMS"
		}]
	},
	{
		"cms": "Kleeja",
		"rules": [{
			"type": "body",
			"keyword": "Powered by Kleeja"
		}]
	},
	{
		"cms": "lepton-cms",
		"rules": [{
			"type": "body",
			"keyword": "Powered by LEPTON CMS"
		}]
	},
	{
		"cms": "MallBuilder",
		"rules": [{
			"type": "body",
			"keyword": "Powered by MallBuilder"
		}]
	},
	{
		"cms": "MediaWiki",
		"rules": [{
			"type": "body",
			"keyword": "Powered by MediaWiki"
		}]
	},
	{
		"cms": "MoMoCMS",
		"rules": [{
			"type": "body",
			"keyword": "Powered BY MoMoCMS"
		}]
	},
	{
		"cms": "OpenCart",
		"rules": [{
			"type": "body",
			"keyword": "Powered By OpenCart"
		}]
	},
	{
		"cms": "opencms",
		"rules": [{
			"type": "body",
			"keyword": "Powered by OpenCms"
		}]
	},
	{
		"cms": "ourphp",
		"rules": [{
			"type": "body",
			"keyword": "Powered by ourphp"
		}]
	},
	{
		"cms": "phpb2b",
		"rules": [{
			"type": "body",
			"keyword": "Powered By PHPB2B"
		}]
	},
	{
		"cms": "PhpCMS",
		"rules": [{
			"type": "body",
			"keyword": "Powered by Phpcms"
		}]
	},
	{
		"cms": "phpdisk",
		"rules": [{
			"type": "body",
			"keyword": "Powered by PHPDisk"
		}]
	},
	{
		"cms": "phpmps",
		"rules": [{
			"type": "body",
			"keyword": "Powered by Phpmps"
		}]
	},
	{
		"cms": "phpok",
		"rules": [{
			"type": "body",
			"keyword": "Powered By phpok.com"
		}]
	},
	{
		"cms": "phpshe",
		"rules": [{
			"type": "body",
			"keyword": "Powered by phpshe"
		}]
	},
	{
		"cms": "phpvod",
		"rules": [{
			"type": "body",
			"keyword": "Powered by PHPVOD"
		}]
	},
	{
		"cms": "海洋CMS",
		"rules": [{
			"type": "body",
			"keyword": "Powered by SeaCms"
		}]
	},
	{
		"cms": "shopbuilder",
		"rules": [{
			"type": "body",
			"keyword": "Powered by ShopBuilder"
		}]
	},
	{
		"cms": "shopnc",
		"rules": [{
			"type": "body",
			"keyword": "Powered by ShopNC"
		}]
	},
	{
		"cms": "ThinkOX",
		"rules": [{
			"type": "body",
			"keyword": "Powered By ThinkOX"
		}]
	},
	{
		"cms": "TurboCMS",
		"rules": [{
			"type": "body",
			"keyword": "Powered by TurboCMS"
		}]
	},
	{
		"cms": "TurboMail",
		"rules": [{
			"type": "body",
			"keyword": "Powered by TurboMail"
		}]
	},
	{
		"cms": "tutucms",
		"rules": [{
			"type": "body",
			"keyword": "Powered by TUTUCMS"
		}]
	},
	{
		"cms": "v5shop",
		"rules": [{
			"type": "body",
			"keyword": "Powered by V5Shop"
		}]
	},
	{
		"cms": "Vicworl",
		"rules": [{
			"type": "body",
			"keyword": "Powered by Vicworl"
		}]
	},
	{
		"cms": "wuzhicms",
		"rules": [{
			"type": "body",
			"keyword": "Powered by wuzhicms"
		}]
	},
	{
		"cms": "Foosun",
		"rules": [{
			"type": "body",
			"keyword": "Powered by www.Foosun.net,Products:Foosun Content Manage system"
		}]
	},
	{
		"cms": "微普外卖点餐系统",
		"rules": [{
			"type": "body",
			"keyword": "Powered By 点餐系统"
		}]
	},
	{
		"cms": "Censura",
		"rules": [{
			"type": "body",
			"keyword": "Powered by: <a href=\"http://www.censura.info"
		}]
	},
	{
		"cms": "Basic-PHP-Events-Lister",
		"rules": [{
			"type": "body",
			"keyword": "Powered by: <a href=\"http://www.mevin.com/\">"
		}]
	},
	{
		"cms": "Amiro-CMS",
		"rules": [{
			"type": "body",
			"keyword": "Powered by: Amiro CMS"
		}]
	},
	{
		"cms": "Arab-Portal",
		"rules": [{
			"type": "body",
			"keyword": "Powered by: Arab"
		}]
	},
	{
		"cms": "Isolsoft-Support-Center",
		"rules": [{
			"type": "body",
			"keyword": "Powered by: Support Center"
		}]
	},
	{
		"cms": "Zotonic",
		"rules": [{
			"type": "body",
			"keyword": "powered by: Zotonic"
		}]
	},
	{
		"cms": "MetInfo",
		"rules": [{
			"type": "body",
			"keyword": "powered_by_metinfo"
		}]
	},
	{
		"cms": "Allomani",
		"rules": [{
			"type": "body",
			"keyword": "Programmed By Allomani"
		}]
	},
	{
		"cms": "Axis-PrintServer",
		"rules": [{
			"type": "body",
			"keyword": "psb_printjobs.gif"
		}]
	},
	{
		"cms": "BoyowCMS",
		"rules": [{
			"type": "body",
			"keyword": "publish by BoyowCMS"
		}]
	},
	{
		"cms": "捷点JCMS",
		"rules": [{
			"type": "body",
			"keyword": "Publish By JCms2010"
		}]
	},
	{
		"cms": "kesionCMS",
		"rules": [{
			"type": "body",
			"keyword": "publish by KesionCMS"
		}]
	},
	{
		"cms": "e-tiller",
		"rules": [{
			"type": "body",
			"keyword": "reader/view_abstract.aspx"
		}]
	},
	{
		"cms": "recaptcha",
		"rules": [{
			"type": "body",
			"keyword": "recaptcha_ajax.js"
		}]
	},
	{
		"cms": "DMXReady-Portfolio-Manager",
		"rules": [{
			"type": "body",
			"keyword": "rememberme_portfoliomanager"
		}]
	},
	{
		"cms": "Sophos Web Appliance",
		"rules": [{
			"type": "body",
			"keyword": "resources/images/sophos_web.ico"
		}]
	},
	{
		"cms": "深信服防火墙类产品",
		"rules": [{
			"type": "body",
			"keyword": "SANGFOR FW"
		}]
	},
	{
		"cms": "SEMcms",
		"rules": [{
			"type": "body",
			"keyword": "sc_mid_c_left_c sc_mid_left_bt"
		}]
	},
	{
		"cms": "逐浪zoomla",
		"rules": [{
			"type": "body",
			"keyword": "script src=\"http://code.zoomla.cn/"
		}]
	},
	{
		"cms": "Energine",
		"rules": [{
			"type": "body",
			"keyword": "scripts/Energine.js"
		}]
	},
	{
		"cms": "SEMcms",
		"rules": [{
			"type": "body",
			"keyword": "semcms PHP"
		}]
	},
	{
		"cms": "H3C公司产品",
		"rules": [{
			"type": "body",
			"keyword": "service@h3c.com"
		}]
	},
	{
		"cms": "青果软件",
		"rules": [{
			"type": "body",
			"keyword": "SetKingoEncypt.jsp"
		}]
	},
	{
		"cms": "shopbuilder",
		"rules": [{
			"type": "body",
			"keyword": "ShopBuilder版权所有"
		}]
	},
	{
		"cms": "BIGACE",
		"rules": [{
			"type": "body",
			"keyword": "Site is running BIGACE"
		}]
	},
	{
		"cms": "cart_engine",
		"rules": [{
			"type": "body",
			"keyword": "skins/_common/jscripts.css"
		}]
	},
	{
		"cms": "Solr",
		"rules": [{
			"type": "body",
			"keyword": "SolrCore Initialization Failures"
		}]
	},
	{
		"cms": "techbridge",
		"rules": [{
			"type": "body",
			"keyword": "Sorry,you need to use IE brower"
		}]
	},
	{
		"cms": "北京金盘鹏图软件",
		"rules": [{
			"type": "body",
			"keyword": "SpeakIntertScarch.aspx"
		}]
	},
	{
		"cms": "3COM NBX",
		"rules": [{
			"type": "body",
			"keyword": "splashTitleIPTelephony"
		}]
	},
	{
		"cms": "splunk",
		"rules": [{
			"type": "body",
			"keyword": "Splunk.util.normalizeBoolean"
		}]
	},
	{
		"cms": "帝友P2P",
		"rules": [{
			"type": "body",
			"keyword": "src=\"/dyweb/dythemes"
		}]
	},
	{
		"cms": "Kloxo-Single-Server",
		"rules": [{
			"type": "body",
			"keyword": "src=\"/img/hypervm-logo.gif"
		}]
	},
	{
		"cms": "贷齐乐p2p",
		"rules": [{
			"type": "body",
			"keyword": "src=\"/js/jPackage"
		}]
	},
	{
		"cms": "78oa",
		"rules": [{
			"type": "body",
			"keyword": "src=\"/module/index.php"
		}]
	},
	{
		"cms": "Acidcat_CMS",
		"rules": [{
			"type": "body",
			"keyword": "Start Acidcat CMS footer information"
		}]
	},
	{
		"cms": "Acidcat CMS",
		"rules": [{
			"type": "body",
			"keyword": "Start Acidcat CMS footer information"
		}]
	},
	{
		"cms": "Kampyle",
		"rules": [{
			"type": "body",
			"keyword": "Start Kampyle Feedback Form Button"
		}]
	},
	{
		"cms": "万网企业云邮箱",
		"rules": [{
			"type": "body",
			"keyword": "static.mxhichina.com/images/favicon.ico"
		}]
	},
	{
		"cms": "Storm",
		"rules": [{
			"type": "body",
			"keyword": "stormtimestr"
		}]
	},
	{
		"cms": "同城多用户商城",
		"rules": [{
			"type": "body",
			"keyword": "style_chaoshi"
		}]
	},
	{
		"cms": "正方教务管理系统",
		"rules": [{
			"type": "body",
			"keyword": "style/base/jw.css"
		}]
	},
	{
		"cms": "正方教务管理系统",
		"rules": [{
			"type": "body",
			"keyword": "style/base/jw.css"
		}]
	},
	{
		"cms": "DiBos",
		"rules": [{
			"type": "body",
			"keyword": "style/bovisnt.css"
		}]
	},
	{
		"cms": "DD-WRT",
		"rules": [{
			"type": "body",
			"keyword": "style/pwc/ddwrt.css"
		}]
	},
	{
		"cms": "Energine",
		"rules": [{
			"type": "body",
			"keyword": "stylesheets/energine.css"
		}]
	},
	{
		"cms": "ASProxy",
		"rules": [{
			"type": "body",
			"keyword": "Surf the web invisibly using ASProxy power"
		}]
	},
	{
		"cms": "华为（HUAWEI）安全设备",
		"rules": [{
			"type": "body",
			"keyword": "sweb-lib/resource/"
		}]
	},
	{
		"cms": "Synology_DiskStation",
		"rules": [{
			"type": "body",
			"keyword": "SYNO.SDS.Session"
		}]
	},
	{
		"cms": "Contao",
		"rules": [{
			"type": "body",
			"keyword": "system/contao.css"
		}]
	},
	{
		"cms": "SiteServer",
		"rules": [{
			"type": "body",
			"keyword": "T_系统首页模板"
		}]
	},
	{
		"cms": "eLitius",
		"rules": [{
			"type": "body",
			"keyword": "target=\"_blank\" title=\"Affiliate"
		}]
	},
	{
		"cms": "Claroline",
		"rules": [{
			"type": "body",
			"keyword": "target=\"_blank\">Claroline</a>"
		}]
	},
	{
		"cms": "eDirectory",
		"rules": [{
			"type": "body",
			"keyword": "target=\"_blank\">eDirectory&trade"
		}]
	},
	{
		"cms": "Help-Desk-Software",
		"rules": [{
			"type": "body",
			"keyword": "target=\"_blank\">freehelpdesk.org"
		}]
	},
	{
		"cms": "FrogCMS",
		"rules": [{
			"type": "body",
			"keyword": "target=\"_blank\">Frog CMS"
		}]
	},
	{
		"cms": "Telerik Sitefinity",
		"rules": [{
			"type": "body",
			"keyword": "Telerik.Web.UI.WebResource.axd"
		}]
	},
	{
		"cms": "beecms",
		"rules": [{
			"type": "body",
			"keyword": "template/default/images/slides.min.jquery.js"
		}]
	},
	{
		"cms": "phpmps",
		"rules": [{
			"type": "body",
			"keyword": "templates/phpmps/style/index.css"
		}]
	},
	{
		"cms": "testlink",
		"rules": [{
			"type": "body",
			"keyword": "testlink_library.js"
		}]
	},
	{
		"cms": "BlueOnyx",
		"rules": [{
			"type": "body",
			"keyword": "Thank you for using the BlueOnyx"
		}]
	},
	{
		"cms": "MVB2000",
		"rules": [{
			"type": "body",
			"keyword": "The Magic Voice Box"
		}]
	},
	{
		"cms": "MVB2000",
		"rules": [{
			"type": "body",
			"keyword": "The Magic Voice Box"
		}]
	},
	{
		"cms": "ApPHP-Calendar",
		"rules": [{
			"type": "body",
			"keyword": "This script was generated by ApPHP Calendar"
		}]
	},
	{
		"cms": "CameraLife",
		"rules": [{
			"type": "body",
			"keyword": "This site is powered by Camera Life"
		}]
	},
	{
		"cms": "Axous",
		"rules": [{
			"type": "body",
			"keyword": "title=\"Axous Shareware Shop"
		}]
	},
	{
		"cms": "Edito-CMS",
		"rules": [{
			"type": "body",
			"keyword": "title=\"CMS\" href=\"http://www.edito.pl/"
		}]
	},
	{
		"cms": "CruxCMS",
		"rules": [{
			"type": "body",
			"keyword": "title=\"CruxCMS\" class=\"blank"
		}]
	},
	{
		"cms": "FestOS",
		"rules": [{
			"type": "body",
			"keyword": "title=\"FestOS"
		}]
	},
	{
		"cms": "Custom-CMS",
		"rules": [{
			"type": "body",
			"keyword": "title=\"Powered by CCMS"
		}]
	},
	{
		"cms": "FreeNAS",
		"rules": [{
			"type": "body",
			"keyword": "title=\"Welcome to FreeNAS"
		}]
	},
	{
		"cms": "b2bbuilder",
		"rules": [{
			"type": "body",
			"keyword": "translateButtonId = \"B2Bbuilder"
		}]
	},
	{
		"cms": "teamportal",
		"rules": [{
			"type": "body",
			"keyword": "TS_expiredurl"
		}]
	},
	{
		"cms": "tutucms",
		"rules": [{
			"type": "body",
			"keyword": "TUTUCMS\""
		}]
	},
	{
		"cms": "科迈RAS系统",
		"rules": [{
			"type": "body",
			"keyword": "type=\"application/npRas"
		}]
	},
	{
		"cms": "08cms",
		"rules": [{
			"type": "body",
			"keyword": "typeof(_08cms)"
		}]
	},
	{
		"cms": "moosefs",
		"rules": [{
			"type": "body",
			"keyword": "under-goal files"
		}]
	},
	{
		"cms": "moosefs",
		"rules": [{
			"type": "body",
			"keyword": "under-goal files"
		}]
	},
	{
		"cms": "asp168欧虎",
		"rules": [{
			"type": "body",
			"keyword": "upload/moban/images/style.css"
		}]
	},
	{
		"cms": "Sophos Web Appliance",
		"rules": [{
			"type": "body",
			"keyword": "url(resources/images/en/login_swa.jpg)"
		}]
	},
	{
		"cms": "单点CRM系统",
		"rules": [{
			"type": "body",
			"keyword": "URL=general/ERP/LOGIN/"
		}]
	},
	{
		"cms": "微普外卖点餐系统",
		"rules": [{
			"type": "body",
			"keyword": "userfiles/shoppics/"
		}]
	},
	{
		"cms": "euse_study",
		"rules": [{
			"type": "body",
			"keyword": "UserInfo/UserFP.aspx"
		}]
	},
	{
		"cms": "蓝凌EIS智慧协同平台",
		"rules": [{
			"type": "body",
			"keyword": "v11_QRcodeBar clr"
		}]
	},
	{
		"cms": "BugTracker.NET",
		"rules": [{
			"type": "body",
			"keyword": "valign=middle><a href=http://ifdefined.com/bugtrackernet.html>"
		}]
	},
	{
		"cms": "ALCASAR",
		"rules": [{
			"type": "body",
			"keyword": "valoriserDiv5"
		}]
	},
	{
		"cms": "BlueQuartz",
		"rules": [{
			"type": "body",
			"keyword": "VALUE=\"Copyright (C) 2000, Cobalt Networks"
		}]
	},
	{
		"cms": "Evo-Cam",
		"rules": [{
			"type": "body",
			"keyword": "value=\"evocam.jar"
		}]
	},
	{
		"cms": "Vicworl",
		"rules": [{
			"type": "body",
			"keyword": "vindex_right_d"
		}]
	},
	{
		"cms": "Avaya-Aura-Utility-Server",
		"rules": [{
			"type": "body",
			"keyword": "vmsTitle\">Avaya Aura&#8482;&nbsp;Utility Server"
		}]
	},
	{
		"cms": "BEA-WebLogic-Server",
		"rules": [{
			"type": "body",
			"keyword": "WebLogic"
		}]
	},
	{
		"cms": "Webmin",
		"rules": [{
			"type": "body",
			"keyword": "Webmin server on"
		}]
	},
	{
		"cms": "Webmin",
		"rules": [{
			"type": "body",
			"keyword": "Webmin server on"
		}]
	},
	{
		"cms": "wecenter",
		"rules": [{
			"type": "body",
			"keyword": "WeCenter"
		}]
	},
	{
		"cms": "FileVista",
		"rules": [{
			"type": "body",
			"keyword": "Welcome to FileVista"
		}]
	},
	{
		"cms": "Advanced-Image-Hosting-Script",
		"rules": [{
			"type": "body",
			"keyword": "Welcome to install AIHS Script"
		}]
	},
	{
		"cms": "juniper_vpn",
		"rules": [{
			"type": "body",
			"keyword": "welcome.cgi?p=logo"
		}]
	},
	{
		"cms": "Astaro-Security-Gateway",
		"rules": [{
			"type": "body",
			"keyword": "wfe/asg/js/app_selector.js?t="
		}]
	},
	{
		"cms": "ezOFFICE",
		"rules": [{
			"type": "body",
			"keyword": "whirRootPath"
		}]
	},
	{
		"cms": "擎天电子政务",
		"rules": [{
			"type": "body",
			"keyword": "window.location = \"homepages/index.aspx"
		}]
	},
	{
		"cms": "trs_wcm",
		"rules": [{
			"type": "body",
			"keyword": "window.location.href = \"/wcm\";"
		}]
	},
	{
		"cms": "Citrix-Metaframe",
		"rules": [{
			"type": "body",
			"keyword": "window.location=\"/Citrix/MetaFrame"
		}]
	},
	{
		"cms": "锐捷应用控制引擎",
		"rules": [{
			"type": "body",
			"keyword": "window.open(\"/login.do\",\"airWin"
		}]
	},
	{
		"cms": "winwebmail",
		"rules": [{
			"type": "body",
			"keyword": "WinWebMail Server"
		}]
	},
	{
		"cms": "WishOA",
		"rules": [{
			"type": "body",
			"keyword": "WishOA_WebPlugin.js"
		}]
	},
	{
		"cms": "WordPress-php",
		"rules": [{
			"type": "body",
			"keyword": "wp-user"
		}]
	},
	{
		"cms": "Google-Talk-Chatback",
		"rules": [{
			"type": "body",
			"keyword": "www.google.com/talk/service/"
		}]
	},
	{
		"cms": "TurboMail",
		"rules": [{
			"type": "body",
			"keyword": "wzcon1 clearfix"
		}]
	},
	{
		"cms": "xheditor",
		"rules": [{
			"type": "body",
			"keyword": "xheditor_lang/zh-cn.js"
		}]
	},
	{
		"cms": "Apache-Wicket",
		"rules": [{
			"type": "body",
			"keyword": "xmlns:wicket="
		}]
	},
	{
		"cms": "bit-service",
		"rules": [{
			"type": "body",
			"keyword": "xmlpzs/webissue.asp"
		}]
	},
	{
		"cms": "OnSSI_Video_Clients",
		"rules": [{
			"type": "body",
			"keyword": "x-value=\"On-Net Surveillance Systems Inc.\""
		}]
	},
	{
		"cms": "全国烟草系统",
		"rules": [{
			"type": "body",
			"keyword": "ycportal/webpublish"
		}]
	},
	{
		"cms": "yidacms",
		"rules": [{
			"type": "body",
			"keyword": "yidacms.css"
		}]
	},
	{
		"cms": "元年财务软件",
		"rules": [{
			"type": "body",
			"keyword": "yuannian.css"
		}]
	},
	{
		"cms": "ZCMS",
		"rules": [{
			"type": "body",
			"keyword": "zcms_skin"
		}]
	},
	{
		"cms": "正方OA",
		"rules": [{
			"type": "body",
			"keyword": "zfoausername"
		}]
	},
	{
		"cms": "智睿软件",
		"rules": [{
			"type": "body",
			"keyword": "Zhirui.js"
		}]
	},
	{
		"cms": "ZoneMinder",
		"rules": [{
			"type": "body",
			"keyword": "ZoneMinder Login"
		}]
	},
	{
		"cms": "信达OA",
		"rules": [{
			"type": "body",
			"keyword": "北京创信达科技有限公司"
		}]
	},
	{
		"cms": "URP教务系统",
		"rules": [{
			"type": "body",
			"keyword": "北京清元优软科技有限公司"
		}]
	},
	{
		"cms": "中国期刊先知网",
		"rules": [{
			"type": "body",
			"keyword": "本系统由<span class=\"STYLE1\" ><a href=\"http://www.firstknow.cn"
		}]
	},
	{
		"cms": "凡科",
		"rules": [{
			"type": "body",
			"keyword": "凡科互联网科技股份有限公司"
		}]
	},
	{
		"cms": "科来RAS",
		"rules": [{
			"type": "body",
			"keyword": "科来软件 版权所有"
		}]
	},
	{
		"cms": "易普拉格科研管理系统",
		"rules": [{
			"type": "body",
			"keyword": "科研管理系统，北京易普拉格科技"
		}]
	},
	{
		"cms": "主机宝",
		"rules": [{
			"type": "body",
			"keyword": "您访问的是主机宝服务器默认页"
		}]
	},
	{
		"cms": "bxemail",
		"rules": [{
			"type": "body",
			"keyword": "请输入正确的电子邮件地址，如：abc@bxemail.com"
		}]
	},
	{
		"cms": "百为路由",
		"rules": [{
			"type": "body",
			"keyword": "提交验证的id必须是ctl_submit"
		}]
	},
	{
		"cms": "天融信TopFlow",
		"rules": [{
			"type": "body",
			"keyword": "天融信TopFlow"
		}]
	},
	{
		"cms": "天融信入侵防御系统TopIDP",
		"rules": [{
			"type": "body",
			"keyword": "天融信入侵防御系统TopIDP"
		}]
	},
	{
		"cms": "沃科网异网同显系统",
		"rules": [{
			"type": "body",
			"keyword": "沃科网"
		}]
	},
	{
		"cms": "武汉弘智科技",
		"rules": [{
			"type": "body",
			"keyword": "研发与技术支持：武汉弘智科技有限公司"
		}]
	},
	{
		"cms": "javashop",
		"rules": [{
			"type": "body",
			"keyword": "易族智汇javashop"
		}]
	},
	{
		"cms": "科迈RAS系统",
		"rules": [{
			"type": "body",
			"keyword": "远程技术支持请求：<a href=\"http://www.comexe.cn"
		}]
	},
	{
		"cms": "中企动力门户CMS",
		"rules": [{
			"type": "body",
			"keyword": "中企动力提供技术支持"
		}]
	},
	{
		"cms": "Coremail",
		"rules": [{
			"type": "title",
			"keyword": "/coremail/common/assets"
		}]
	},
	{
		"cms": "171cms",
		"rules": [{
			"type": "title",
			"keyword": "171cms"
		}]
	},
	{
		"cms": "78oa",
		"rules": [{
			"type": "title",
			"keyword": "78oa"
		}]
	},
	{
		"cms": "网动云视讯平台",
		"rules": [{
			"type": "title",
			"keyword": "Acenter"
		}]
	},
	{
		"cms": "Adiscon_LogAnalyzer",
		"rules": [{
			"type": "title",
			"keyword": "Adiscon LogAnalyzer"
		}]
	},
	{
		"cms": "AirTiesRouter",
		"rules": [{
			"type": "title",
			"keyword": "Airties"
		}]
	},
	{
		"cms": "H3C AM8000",
		"rules": [{
			"type": "title",
			"keyword": "AM8000"
		}]
	},
	{
		"cms": "AnyGate",
		"rules": [{
			"type": "title",
			"keyword": "AnyGate"
		}]
	},
	{
		"cms": "AP-Router",
		"rules": [{
			"type": "title",
			"keyword": "AP Router New Generation"
		}]
	},
	{
		"cms": "Apache-Archiva",
		"rules": [{
			"type": "title",
			"keyword": "Apache Archiva"
		}]
	},
	{
		"cms": "AVCON6",
		"rules": [{
			"type": "title",
			"keyword": "AVCON6系统管理平台"
		}]
	},
	{
		"cms": "Axis-Network-Camera",
		"rules": [{
			"type": "title",
			"keyword": "AXIS Video Server"
		}]
	},
	{
		"cms": "bacula-web",
		"rules": [{
			"type": "title",
			"keyword": "Bacula Web"
		}]
	},
	{
		"cms": "bacula-web",
		"rules": [{
			"type": "title",
			"keyword": "Bacula-Web"
		}]
	},
	{
		"cms": "bacula-web",
		"rules": [{
			"type": "title",
			"keyword": "bacula-web"
		}]
	},
	{
		"cms": "baocms",
		"rules": [{
			"type": "title",
			"keyword": "baocms"
		}]
	},
	{
		"cms": "Barracuda-Spam-Firewall",
		"rules": [{
			"type": "title",
			"keyword": "Barracuda Spam & Virus Firewall: Welcome"
		}]
	},
	{
		"cms": "BigDump",
		"rules": [{
			"type": "title",
			"keyword": "BigDump"
		}]
	},
	{
		"cms": "Biromsoft-WebCam",
		"rules": [{
			"type": "title",
			"keyword": "Biromsoft WebCam"
		}]
	},
	{
		"cms": "BlueNet-Video",
		"rules": [{
			"type": "title",
			"keyword": "BlueNet Video Viewer Version"
		}]
	},
	{
		"cms": "BugFree",
		"rules": [{
			"type": "title",
			"keyword": "BugFree"
		}]
	},
	{
		"cms": "CalendarScript",
		"rules": [{
			"type": "title",
			"keyword": "Calendar Administration : Login"
		}]
	},
	{
		"cms": "CDR-Stats",
		"rules": [{
			"type": "title",
			"keyword": "CDR-Stats | Customer Interface"
		}]
	},
	{
		"cms": "Centreon",
		"rules": [{
			"type": "title",
			"keyword": "Centreon - IT & Network Monitoring"
		}]
	},
	{
		"cms": "CGI:IRC",
		"rules": [{
			"type": "title",
			"keyword": "CGI:IRC Login"
		}]
	},
	{
		"cms": "AChecker Web accessibility evaluation tool",
		"rules": [{
			"type": "title",
			"keyword": "Checker : Web Accessibility Checker"
		}]
	},
	{
		"cms": "Cisco_Cable_Modem",
		"rules": [{
			"type": "title",
			"keyword": "Cisco Cable Modem"
		}]
	},
	{
		"cms": "Cisco-VPN-3000-Concentrator",
		"rules": [{
			"type": "title",
			"keyword": "Cisco Systems, Inc. VPN 3000 Concentrator"
		}]
	},
	{
		"cms": "cisco UCM",
		"rules": [{
			"type": "title",
			"keyword": "Cisco Unified"
		}]
	},
	{
		"cms": "Cogent-DataHub",
		"rules": [{
			"type": "title",
			"keyword": "Cogent DataHub WebView"
		}]
	},
	{
		"cms": "cPassMan",
		"rules": [{
			"type": "title",
			"keyword": "Collaborative Passwords Manager"
		}]
	},
	{
		"cms": "Coremail",
		"rules": [{
			"type": "title",
			"keyword": "Coremail邮件系统"
		}]
	},
	{
		"cms": "DVWA",
		"rules": [{
			"type": "title",
			"keyword": "Damn Vulnerable Web App (DVWA) - Login"
		}]
	},
	{
		"cms": "D-Link-Network-Camera",
		"rules": [{
			"type": "title",
			"keyword": "DCS-5300"
		}]
	},
	{
		"cms": "Dell-Printer",
		"rules": [{
			"type": "title",
			"keyword": "Dell Laser Printer"
		}]
	},
	{
		"cms": "Dell OpenManage Switch Administrator",
		"rules": [{
			"type": "title",
			"keyword": "Dell OpenManage Switch Administrator"
		}]
	},
	{
		"cms": "DiBos",
		"rules": [{
			"type": "title",
			"keyword": "DiBos - Login"
		}]
	},
	{
		"cms": "D-Link_VoIP_Wireless_Router",
		"rules": [{
			"type": "title",
			"keyword": "D-Link VoIP Wireless Router"
		}]
	},
	{
		"cms": "Dorado",
		"rules": [{
			"type": "title",
			"keyword": "Dorado Login Page"
		}]
	},
	{
		"cms": "DORG",
		"rules": [{
			"type": "title",
			"keyword": "DORG - "
		}]
	},
	{
		"cms": "dtcms",
		"rules": [{
			"type": "title",
			"keyword": "dtcms"
		}]
	},
	{
		"cms": "DVR camera",
		"rules": [{
			"type": "title",
			"keyword": "DVR WebClient"
		}]
	},
	{
		"cms": "DVR-WebClient",
		"rules": [{
			"type": "title",
			"keyword": "DVR-WebClient"
		}]
	},
	{
		"cms": "eadmin",
		"rules": [{
			"type": "title",
			"keyword": "eadmin"
		}]
	},
	{
		"cms": "eBuilding-Network-Controller",
		"rules": [{
			"type": "title",
			"keyword": "eBuilding Web"
		}]
	},
	{
		"cms": "EDIMAX",
		"rules": [{
			"type": "title",
			"keyword": "EDIMAX Technology"
		}]
	},
	{
		"cms": "EdmWebVideo",
		"rules": [{
			"type": "title",
			"keyword": "EdmWebVideo"
		}]
	},
	{
		"cms": "EduSoho开源网络课堂",
		"rules": [{
			"type": "title",
			"keyword": "edusoho"
		}]
	},
	{
		"cms": "edvr",
		"rules": [{
			"type": "title",
			"keyword": "edvs/edvr"
		}]
	},
	{
		"cms": "Entrans",
		"rules": [{
			"type": "title",
			"keyword": "Entrans"
		}]
	},
	{
		"cms": "H3C ER2100n",
		"rules": [{
			"type": "title",
			"keyword": "ER2100n系统管理"
		}]
	},
	{
		"cms": "H3C ER2100V2",
		"rules": [{
			"type": "title",
			"keyword": "ER2100V2系统管理"
		}]
	},
	{
		"cms": "H3C ER2100",
		"rules": [{
			"type": "title",
			"keyword": "ER2100系统管理"
		}]
	},
	{
		"cms": "H3C ER3100",
		"rules": [{
			"type": "title",
			"keyword": "ER3100系统管理"
		}]
	},
	{
		"cms": "H3C ER3108GW",
		"rules": [{
			"type": "title",
			"keyword": "ER3108GW系统管理"
		}]
	},
	{
		"cms": "H3C ER3108G",
		"rules": [{
			"type": "title",
			"keyword": "ER3108G系统管理"
		}]
	},
	{
		"cms": "H3C ER3200",
		"rules": [{
			"type": "title",
			"keyword": "ER3200系统管理"
		}]
	},
	{
		"cms": "H3C ER3260G2",
		"rules": [{
			"type": "title",
			"keyword": "ER3260G2系统管理"
		}]
	},
	{
		"cms": "H3C ER3260",
		"rules": [{
			"type": "title",
			"keyword": "ER3260系统管理"
		}]
	},
	{
		"cms": "H3C ER5100",
		"rules": [{
			"type": "title",
			"keyword": "ER5100系统管理"
		}]
	},
	{
		"cms": "H3C ER5200G2",
		"rules": [{
			"type": "title",
			"keyword": "ER5200G2系统管理"
		}]
	},
	{
		"cms": "H3C ER5200",
		"rules": [{
			"type": "title",
			"keyword": "ER5200系统管理"
		}]
	},
	{
		"cms": "H3C ER6300G2",
		"rules": [{
			"type": "title",
			"keyword": "ER6300G2系统管理"
		}]
	},
	{
		"cms": "H3C ER6300",
		"rules": [{
			"type": "title",
			"keyword": "ER6300系统管理"
		}]
	},
	{
		"cms": "H3C ER8300G2",
		"rules": [{
			"type": "title",
			"keyword": "ER8300G2系统管理"
		}]
	},
	{
		"cms": "H3C ER8300",
		"rules": [{
			"type": "title",
			"keyword": "ER8300系统管理"
		}]
	},
	{
		"cms": "yongyoufe",
		"rules": [{
			"type": "title",
			"keyword": "FE协作"
		}]
	},
	{
		"cms": "File-Upload-Manager",
		"rules": [{
			"type": "title",
			"keyword": "File Upload Manager"
		}]
	},
	{
		"cms": "Fortinet Firewall",
		"rules": [{
			"type": "title",
			"keyword": "Firewall Notification"
		}]
	},
	{
		"cms": "Forest-Blog",
		"rules": [{
			"type": "title",
			"keyword": "Forest Blog"
		}]
	},
	{
		"cms": "DnP-Firewall",
		"rules": [{
			"type": "title",
			"keyword": "Forum Gateway - Powered by DnP Firewall"
		}]
	},
	{
		"cms": "FreeboxOS",
		"rules": [{
			"type": "title",
			"keyword": "Freebox OS"
		}]
	},
	{
		"cms": "Gallarific",
		"rules": [{
			"type": "title",
			"keyword": "Gallarific > Sign in"
		}]
	},
	{
		"cms": "Gallery",
		"rules": [{
			"type": "title",
			"keyword": "Gallery 3 Installer"
		}]
	},
	{
		"cms": "GateQuest-PHP-Site-Recommender",
		"rules": [{
			"type": "title",
			"keyword": "GateQuest"
		}]
	},
	{
		"cms": "GenieATM",
		"rules": [{
			"type": "title",
			"keyword": "GenieATM"
		}]
	},
	{
		"cms": "GenOHM-SCADA",
		"rules": [{
			"type": "title",
			"keyword": "GenOHM Scada Launcher"
		}]
	},
	{
		"cms": "Gossamer-Forum",
		"rules": [{
			"type": "title",
			"keyword": "Gossamer Forum"
		}]
	},
	{
		"cms": "GpsGate-Server",
		"rules": [{
			"type": "title",
			"keyword": "GpsGate Server - "
		}]
	},
	{
		"cms": "GPSweb",
		"rules": [{
			"type": "title",
			"keyword": "GPSweb"
		}]
	},
	{
		"cms": "Honeywell IP-Camera",
		"rules": [{
			"type": "title",
			"keyword": "Honeywell IP-Camera"
		}]
	},
	{
		"cms": "honeywell NetAXS",
		"rules": [{
			"type": "title",
			"keyword": "Honeywell NetAXS"
		}]
	},
	{
		"cms": "iLO",
		"rules": [{
			"type": "title",
			"keyword": "HP Integrated Lights-Out"
		}]
	},
	{
		"cms": "HP-OfficeJet-Printer",
		"rules": [{
			"type": "title",
			"keyword": "HP Officejet"
		}]
	},
	{
		"cms": "HP-StorageWorks-Library",
		"rules": [{
			"type": "title",
			"keyword": "HP StorageWorks"
		}]
	},
	{
		"cms": "Huawei B683",
		"rules": [{
			"type": "title",
			"keyword": "Huawei B683"
		}]
	},
	{
		"cms": "Huawei B683V",
		"rules": [{
			"type": "title",
			"keyword": "Huawei B683V"
		}]
	},
	{
		"cms": "HUAWEI CSP",
		"rules": [{
			"type": "title",
			"keyword": "HUAWEI CSP"
		}]
	},
	{
		"cms": "HUAWEI ESPACE 7910",
		"rules": [{
			"type": "title",
			"keyword": "HUAWEI ESPACE 7910"
		}]
	},
	{
		"cms": "Huawei HG520 ADSL2+ Router",
		"rules": [{
			"type": "title",
			"keyword": "Huawei HG520"
		}]
	},
	{
		"cms": "Huawei HG630",
		"rules": [{
			"type": "title",
			"keyword": "Huawei HG630"
		}]
	},
	{
		"cms": "HUAWEI Inner Web",
		"rules": [{
			"type": "title",
			"keyword": "HUAWEI Inner Web"
		}]
	},
	{
		"cms": "华为 MCU",
		"rules": [{
			"type": "title",
			"keyword": "huawei MCU"
		}]
	},
	{
		"cms": "华为 NetOpen",
		"rules": [{
			"type": "title",
			"keyword": "Huawei NetOpen System"
		}]
	},
	{
		"cms": "Kloxo-Single-Server",
		"rules": [{
			"type": "title",
			"keyword": "HyperVM"
		}]
	},
	{
		"cms": "i-Gallery",
		"rules": [{
			"type": "title",
			"keyword": "i-Gallery"
		}]
	},
	{
		"cms": "Lotus",
		"rules": [{
			"type": "title",
			"keyword": "IBM Lotus iNotes Login"
		}]
	},
	{
		"cms": "Lotus",
		"rules": [{
			"type": "title",
			"keyword": "IBM Lotus iNotes Login"
		}]
	},
	{
		"cms": "H3C ICG 1000",
		"rules": [{
			"type": "title",
			"keyword": "ICG 1000系统管理"
		}]
	},
	{
		"cms": "H3C ICG1000",
		"rules": [{
			"type": "title",
			"keyword": "ICG1000系统管理"
		}]
	},
	{
		"cms": "I-O-DATA-Router",
		"rules": [{
			"type": "title",
			"keyword": "I-O DATA Wireless Broadband Router"
		}]
	},
	{
		"cms": "iGENUS邮件系统",
		"rules": [{
			"type": "title",
			"keyword": "iGENUS webmail"
		}]
	},
	{
		"cms": "infoglue",
		"rules": [{
			"type": "title",
			"keyword": "infoglue"
		}]
	},
	{
		"cms": "IQeye-Netcam",
		"rules": [{
			"type": "title",
			"keyword": "IQEYE: Live Images"
		}]
	},
	{
		"cms": "iTop",
		"rules": [{
			"type": "title",
			"keyword": "iTop Login"
		}]
	},
	{
		"cms": "jieqi cms",
		"rules": [{
			"type": "title",
			"keyword": "jieqi cms"
		}]
	},
	{
		"cms": "Kibana",
		"rules": [{
			"type": "title",
			"keyword": "Kibana"
		}]
	},
	{
		"cms": "kingcms",
		"rules": [{
			"type": "title",
			"keyword": "kingcms"
		}]
	},
	{
		"cms": "青果软件",
		"rules": [{
			"type": "title",
			"keyword": "KINGOSOFT"
		}]
	},
	{
		"cms": "wdcp管理系统",
		"rules": [{
			"type": "title",
			"keyword": "lanmp_wdcp 安装成功"
		}]
	},
	{
		"cms": "LANMP一键安装包",
		"rules": [{
			"type": "title",
			"keyword": "LANMP一键安装包"
		}]
	},
	{
		"cms": "网御上网行为管理系统",
		"rules": [{
			"type": "title",
			"keyword": "Leadsec ACM"
		}]
	},
	{
		"cms": "Linksys_SPA_Configuration ",
		"rules": [{
			"type": "title",
			"keyword": "Linksys SPA Configuration"
		}]
	},
	{
		"cms": "BlueOnyx",
		"rules": [{
			"type": "title",
			"keyword": "Login - BlueOnyx"
		}]
	},
	{
		"cms": "BlueQuartz",
		"rules": [{
			"type": "title",
			"keyword": "Login - BlueQuartz"
		}]
	},
	{
		"cms": "eXtplorer",
		"rules": [{
			"type": "title",
			"keyword": "Login - eXtplorer"
		}]
	},
	{
		"cms": "Collabtive",
		"rules": [{
			"type": "title",
			"keyword": "Login @ Collabtive"
		}]
	},
	{
		"cms": "Webmin",
		"rules": [{
			"type": "title",
			"keyword": "Login to Webmin"
		}]
	},
	{
		"cms": "Webmin",
		"rules": [{
			"type": "title",
			"keyword": "Login to Webmin"
		}]
	},
	{
		"cms": "LuManager",
		"rules": [{
			"type": "title",
			"keyword": "LuManager"
		}]
	},
	{
		"cms": "Macrec_DVR",
		"rules": [{
			"type": "title",
			"keyword": "Macrec DVR"
		}]
	},
	{
		"cms": "Mercurial",
		"rules": [{
			"type": "title",
			"keyword": "Mercurial repositories index"
		}]
	},
	{
		"cms": "Symantec Messaging Gateway",
		"rules": [{
			"type": "title",
			"keyword": "Messaging Gateway"
		}]
	},
	{
		"cms": "Oracle_OPERA",
		"rules": [{
			"type": "title",
			"keyword": "MICROS Systems Inc., OPERA"
		}]
	},
	{
		"cms": "ZTE_MiFi_UNE",
		"rules": [{
			"type": "title",
			"keyword": "MiFi UNE 4G LTE"
		}]
	},
	{
		"cms": "ZTE_MiFi_UNE",
		"rules": [{
			"type": "title",
			"keyword": "MiFi UNE 4G LTE"
		}]
	},
	{
		"cms": "Mixcall座席管理中心",
		"rules": [{
			"type": "title",
			"keyword": "Mixcall座席管理中心"
		}]
	},
	{
		"cms": "Motorola_SBG900",
		"rules": [{
			"type": "title",
			"keyword": "Motorola SBG900"
		}]
	},
	{
		"cms": "MRTG",
		"rules": [{
			"type": "title",
			"keyword": "MRTG Index Page"
		}]
	},
	{
		"cms": "MVB2000",
		"rules": [{
			"type": "title",
			"keyword": "MVB2000"
		}]
	},
	{
		"cms": "MVB2000",
		"rules": [{
			"type": "title",
			"keyword": "MVB2000"
		}]
	},
	{
		"cms": "mymps",
		"rules": [{
			"type": "title",
			"keyword": "mymps"
		}]
	},
	{
		"cms": "3COM NBX",
		"rules": [{
			"type": "title",
			"keyword": "NBX NetSet"
		}]
	},
	{
		"cms": "NETSurveillance",
		"rules": [{
			"type": "title",
			"keyword": "NETSurveillance"
		}]
	},
	{
		"cms": "ipTIME-Router",
		"rules": [{
			"type": "title",
			"keyword": "networks - ipTIME"
		}]
	},
	{
		"cms": "NOALYSS",
		"rules": [{
			"type": "title",
			"keyword": "NOALYSS"
		}]
	},
	{
		"cms": "绿盟下一代防火墙",
		"rules": [{
			"type": "title",
			"keyword": "NSFOCUS NF"
		}]
	},
	{
		"cms": "soffice",
		"rules": [{
			"type": "title",
			"keyword": "OA办公管理平台"
		}]
	},
	{
		"cms": "OBSERVA telcom",
		"rules": [{
			"type": "title",
			"keyword": "OBSERVA"
		}]
	},
	{
		"cms": "OnSSI_Video_Clients",
		"rules": [{
			"type": "title",
			"keyword": "OnSSI Video Clients"
		}]
	},
	{
		"cms": "openEAP",
		"rules": [{
			"type": "title",
			"keyword": "openEAP_统一登录门户"
		}]
	},
	{
		"cms": "OpenMas",
		"rules": [{
			"type": "title",
			"keyword": "OpenMas"
		}]
	},
	{
		"cms": "OpenMas",
		"rules": [{
			"type": "title",
			"keyword": "OpenMas"
		}]
	},
	{
		"cms": "汉柏安全网关",
		"rules": [{
			"type": "title",
			"keyword": "OPZOON - "
		}]
	},
	{
		"cms": "panabit智能网关",
		"rules": [{
			"type": "title",
			"keyword": "panabit"
		}]
	},
	{
		"cms": "phpems考试系统",
		"rules": [{
			"type": "title",
			"keyword": "phpems"
		}]
	},
	{
		"cms": "phpmoadmin",
		"rules": [{
			"type": "title",
			"keyword": "phpmoadmin"
		}]
	},
	{
		"cms": "phpok",
		"rules": [{
			"type": "title",
			"keyword": "phpok"
		}]
	},
	{
		"cms": "PineApp",
		"rules": [{
			"type": "title",
			"keyword": "PineApp WebAccess - Login"
		}]
	},
	{
		"cms": "DuomiCms",
		"rules": [{
			"type": "title",
			"keyword": "Power by DuomiCms"
		}]
	},
	{
		"cms": "TCCMS",
		"rules": [{
			"type": "title",
			"keyword": "Power By TCCMS"
		}]
	},
	{
		"cms": "ASPCMS",
		"rules": [{
			"type": "title",
			"keyword": "Powered by ASPCMS"
		}]
	},
	{
		"cms": "DedeCMS",
		"rules": [{
			"type": "title",
			"keyword": "Powered by DedeCms\" "
		}]
	},
	{
		"cms": "地平线CMS",
		"rules": [{
			"type": "title",
			"keyword": "Powered by deep soon"
		}]
	},
	{
		"cms": "discuz",
		"rules": [{
			"type": "title",
			"keyword": "Powered by Discuz"
		}]
	},
	{
		"cms": "帝国EmpireCMS",
		"rules": [{
			"type": "title",
			"keyword": "Powered by EmpireCMS"
		}]
	},
	{
		"cms": "ESPCMS",
		"rules": [{
			"type": "title",
			"keyword": "Powered by ESPCMS"
		}]
	},
	{
		"cms": "JEECMS",
		"rules": [{
			"type": "title",
			"keyword": "Powered by JEECMS"
		}]
	},
	{
		"cms": "MetInfo",
		"rules": [{
			"type": "title",
			"keyword": "Powered by MetInfo"
		}]
	},
	{
		"cms": "Npoint",
		"rules": [{
			"type": "title",
			"keyword": "Powered by Npoint"
		}]
	},
	{
		"cms": "phpwind",
		"rules": [{
			"type": "title",
			"keyword": "Powered by phpwind"
		}]
	},
	{
		"cms": "sdcms",
		"rules": [{
			"type": "title",
			"keyword": "powered by sdcms"
		}]
	},
	{
		"cms": "SiteServer",
		"rules": [{
			"type": "title",
			"keyword": "Powered by SiteServer CMS"
		}]
	},
	{
		"cms": "PublicCMS",
		"rules": [{
			"type": "title",
			"keyword": "publiccms"
		}]
	},
	{
		"cms": "Puppet_Node_Manager",
		"rules": [{
			"type": "title",
			"keyword": "Puppet Node Manager"
		}]
	},
	{
		"cms": "OrientDB",
		"rules": [{
			"type": "title",
			"keyword": "Redirecting to OrientDB"
		}]
	},
	{
		"cms": "BMC-Remedy",
		"rules": [{
			"type": "title",
			"keyword": "Remedy Mid Tier"
		}]
	},
	{
		"cms": "RG-PowerCache内容加速系统",
		"rules": [{
			"type": "title",
			"keyword": "RG-PowerCache"
		}]
	},
	{
		"cms": "richmail",
		"rules": [{
			"type": "title",
			"keyword": "Richmail"
		}]
	},
	{
		"cms": "Ruckus",
		"rules": [{
			"type": "title",
			"keyword": "Ruckus Wireless Admin"
		}]
	},
	{
		"cms": "Samsung DVR",
		"rules": [{
			"type": "title",
			"keyword": "Samsung DVR"
		}]
	},
	{
		"cms": "Scientific-Atlanta_Cable_Modem",
		"rules": [{
			"type": "title",
			"keyword": "Scientific-Atlanta Cable Modem"
		}]
	},
	{
		"cms": "Scientific-Atlanta_Cable_Modem",
		"rules": [{
			"type": "title",
			"keyword": "Scientific-Atlanta Cable Modem"
		}]
	},
	{
		"cms": "海洋CMS",
		"rules": [{
			"type": "title",
			"keyword": "seacms"
		}]
	},
	{
		"cms": "网神防火墙",
		"rules": [{
			"type": "title",
			"keyword": "secgate 3600"
		}]
	},
	{
		"cms": "SHOUTcast",
		"rules": [{
			"type": "title",
			"keyword": "SHOUTcast Administrator"
		}]
	},
	{
		"cms": "SIEMENS IP Cameras",
		"rules": [{
			"type": "title",
			"keyword": "SIEMENS IP Camera"
		}]
	},
	{
		"cms": "SLTM32_Configuration",
		"rules": [{
			"type": "title",
			"keyword": "SLTM32 Web Configuration Pages "
		}]
	},
	{
		"cms": "soeasy网站集群系统",
		"rules": [{
			"type": "title",
			"keyword": "SoEasy网站集群"
		}]
	},
	{
		"cms": "Locus_SolarNOC",
		"rules": [{
			"type": "title",
			"keyword": "SolarNOC - Login"
		}]
	},
	{
		"cms": "Solr",
		"rules": [{
			"type": "title",
			"keyword": "Solr Admin"
		}]
	},
	{
		"cms": "sony摄像头",
		"rules": [{
			"type": "title",
			"keyword": "Sony Network Camera"
		}]
	},
	{
		"cms": "Sophos_Web_Appliance",
		"rules": [{
			"type": "title",
			"keyword": "Sophos Web Appliance"
		}]
	},
	{
		"cms": "Sophos Web Appliance",
		"rules": [{
			"type": "title",
			"keyword": "Sophos Web Appliance"
		}]
	},
	{
		"cms": "Spammark邮件信息安全网关",
		"rules": [{
			"type": "title",
			"keyword": "Spammark邮件信息安全网关"
		}]
	},
	{
		"cms": "Spark_Master",
		"rules": [{
			"type": "title",
			"keyword": "Spark Master at"
		}]
	},
	{
		"cms": "Spark_Worker",
		"rules": [{
			"type": "title",
			"keyword": "Spark Worker at"
		}]
	},
	{
		"cms": "srun3000计费认证系统",
		"rules": [{
			"type": "title",
			"keyword": "srun3000"
		}]
	},
	{
		"cms": "Storm",
		"rules": [{
			"type": "title",
			"keyword": "Storm UI"
		}]
	},
	{
		"cms": "Synology_DiskStation",
		"rules": [{
			"type": "title",
			"keyword": "Synology DiskStation"
		}]
	},
	{
		"cms": "ThinkOX",
		"rules": [{
			"type": "title",
			"keyword": "ThinkOX"
		}]
	},
	{
		"cms": "泰信TMailer邮件系统",
		"rules": [{
			"type": "title",
			"keyword": "Tmailer"
		}]
	},
	{
		"cms": "arrisi_Touchstone",
		"rules": [{
			"type": "title",
			"keyword": "Touchstone Status"
		}]
	},
	{
		"cms": "TurboMail",
		"rules": [{
			"type": "title",
			"keyword": "TurboMail邮件系统"
		}]
	},
	{
		"cms": "UcSTAR",
		"rules": [{
			"type": "title",
			"keyword": "UcSTAR 管理控制台"
		}]
	},
	{
		"cms": "UPUPW",
		"rules": [{
			"type": "title",
			"keyword": "UPUPW环境集成包"
		}]
	},
	{
		"cms": "URP教务系统",
		"rules": [{
			"type": "title",
			"keyword": "URP 综合教务系统"
		}]
	},
	{
		"cms": "v5shop",
		"rules": [{
			"type": "title",
			"keyword": "v5shop"
		}]
	},
	{
		"cms": "Verizon_Router",
		"rules": [{
			"type": "title",
			"keyword": "Verizon Router"
		}]
	},
	{
		"cms": "VideoIQ Camera",
		"rules": [{
			"type": "title",
			"keyword": "VideoIQ Camera Login"
		}]
	},
	{
		"cms": "VisualSVN",
		"rules": [{
			"type": "title",
			"keyword": "VisualSVN Server"
		}]
	},
	{
		"cms": "VOS3000",
		"rules": [{
			"type": "title",
			"keyword": "VOS3000"
		}]
	},
	{
		"cms": "VZPP Plesk",
		"rules": [{
			"type": "title",
			"keyword": "VZPP Plesk "
		}]
	},
	{
		"cms": "wamp",
		"rules": [{
			"type": "title",
			"keyword": "WAMPSERVER"
		}]
	},
	{
		"cms": "ezOFFICE",
		"rules": [{
			"type": "title",
			"keyword": "Wanhu ezOFFICE"
		}]
	},
	{
		"cms": "wdcp管理系统",
		"rules": [{
			"type": "title",
			"keyword": "wdcp服务器"
		}]
	},
	{
		"cms": "wdcp",
		"rules": [{
			"type": "title",
			"keyword": "wdcp服务器"
		}]
	},
	{
		"cms": "WDlinux",
		"rules": [{
			"type": "title",
			"keyword": "wdOS"
		}]
	},
	{
		"cms": "ACTi",
		"rules": [{
			"type": "title",
			"keyword": "Web Configurator"
		}]
	},
	{
		"cms": "FortiGuard",
		"rules": [{
			"type": "title",
			"keyword": "Web Filter Block Override"
		}]
	},
	{
		"cms": "bacula-web",
		"rules": [{
			"type": "title",
			"keyword": "Webacula"
		}]
	},
	{
		"cms": "Wimax_CPE",
		"rules": [{
			"type": "title",
			"keyword": "Wimax CPE Configuration"
		}]
	},
	{
		"cms": "Wimax_CPE",
		"rules": [{
			"type": "title",
			"keyword": "Wimax CPE Configuration"
		}]
	},
	{
		"cms": "winwebmail",
		"rules": [{
			"type": "title",
			"keyword": "winwebmail"
		}]
	},
	{
		"cms": "xfinity",
		"rules": [{
			"type": "title",
			"keyword": "Xfinity"
		}]
	},
	{
		"cms": "nvdvr",
		"rules": [{
			"type": "title",
			"keyword": "XWebPlay"
		}]
	},
	{
		"cms": "ZCMS",
		"rules": [{
			"type": "title",
			"keyword": "ZCMS泽元内容管理"
		}]
	},
	{
		"cms": "中控智慧时间安全管理平台",
		"rules": [{
			"type": "title",
			"keyword": "ZKECO 时间&安全管理平台"
		}]
	},
	{
		"cms": "埃森诺网络服务质量检测系统",
		"rules": [{
			"type": "title",
			"keyword": "埃森诺网络服务质量检测系统 "
		}]
	},
	{
		"cms": "bxemail",
		"rules": [{
			"type": "title",
			"keyword": "百姓邮局"
		}]
	},
	{
		"cms": "bxemail",
		"rules": [{
			"type": "title",
			"keyword": "百讯安全邮件系统"
		}]
	},
	{
		"cms": "畅捷通",
		"rules": [{
			"type": "title",
			"keyword": "畅捷通"
		}]
	},
	{
		"cms": "大米CMS",
		"rules": [{
			"type": "title",
			"keyword": "大米CMS-"
		}]
	},
	{
		"cms": "汉码软件",
		"rules": [{
			"type": "title",
			"keyword": "汉码软件"
		}]
	},
	{
		"cms": "护卫神网站安全系统",
		"rules": [{
			"type": "title",
			"keyword": "护卫神.网站安全系统"
		}]
	},
	{
		"cms": "护卫神主机管理",
		"rules": [{
			"type": "title",
			"keyword": "护卫神·主机管理系统"
		}]
	},
	{
		"cms": "金和协同管理平台",
		"rules": [{
			"type": "title",
			"keyword": "金和协同管理平台"
		}]
	},
	{
		"cms": "金龙卡金融化一卡通网站查询子系统",
		"rules": [{
			"type": "title",
			"keyword": "金龙卡金融化一卡通网站查询子系统"
		}]
	},
	{
		"cms": "科来RAS",
		"rules": [{
			"type": "title",
			"keyword": "科来网络回溯"
		}]
	},
	{
		"cms": "科迈RAS系统",
		"rules": [{
			"type": "title",
			"keyword": "科迈RAS"
		}]
	},
	{
		"cms": "单点CRM系统",
		"rules": [{
			"type": "title",
			"keyword": "客户关系管理-CRM"
		}]
	},
	{
		"cms": "浪潮政务系统",
		"rules": [{
			"type": "title",
			"keyword": "浪潮政务"
		}]
	},
	{
		"cms": "任我行CRM",
		"rules": [{
			"type": "title",
			"keyword": "任我行CRM"
		}]
	},
	{
		"cms": "锐捷应用控制引擎",
		"rules": [{
			"type": "title",
			"keyword": "锐捷应用控制引擎"
		}]
	},
	{
		"cms": "瑞友天翼_应用虚拟化系统 ",
		"rules": [{
			"type": "title",
			"keyword": "瑞友天翼－应用虚拟化系统"
		}]
	},
	{
		"cms": "Apabi数字资源平台",
		"rules": [{
			"type": "title",
			"keyword": "数字资源平台"
		}]
	},
	{
		"cms": "ACSNO网络探针",
		"rules": [{
			"type": "title",
			"keyword": "探针管理与测试系统-登录界面"
		}]
	},
	{
		"cms": "天融信 TopAD",
		"rules": [{
			"type": "title",
			"keyword": "天融信 TopAD"
		}]
	},
	{
		"cms": "天融信ADS管理平台",
		"rules": [{
			"type": "title",
			"keyword": "天融信ADS管理平台"
		}]
	},
	{
		"cms": "天融信Web应用安全防护系统",
		"rules": [{
			"type": "title",
			"keyword": "天融信Web应用安全防护系统"
		}]
	},
	{
		"cms": "天融信WEB应用防火墙",
		"rules": [{
			"type": "title",
			"keyword": "天融信WEB应用防火墙"
		}]
	},
	{
		"cms": "天融信安全管理系统",
		"rules": [{
			"type": "title",
			"keyword": "天融信安全管理"
		}]
	},
	{
		"cms": "天融信脆弱性扫描与管理系统",
		"rules": [{
			"type": "title",
			"keyword": "天融信脆弱性扫描与管理系统"
		}]
	},
	{
		"cms": "天融信日志收集与分析系统",
		"rules": [{
			"type": "title",
			"keyword": "天融信日志收集与分析系统"
		}]
	},
	{
		"cms": "天融信入侵检测系统TopSentry",
		"rules": [{
			"type": "title",
			"keyword": "天融信入侵检测系统TopSentry"
		}]
	},
	{
		"cms": "天融信网络卫士过滤网关",
		"rules": [{
			"type": "title",
			"keyword": "天融信网络卫士过滤网关"
		}]
	},
	{
		"cms": "天融信网站监测与自动修复系统",
		"rules": [{
			"type": "title",
			"keyword": "天融信网站监测与自动修复系统"
		}]
	},
	{
		"cms": "天融信异常流量管理与抗拒绝服务系统",
		"rules": [{
			"type": "title",
			"keyword": "天融信异常流量管理与抗拒绝服务系统"
		}]
	},
	{
		"cms": "ezOFFICE",
		"rules": [{
			"type": "title",
			"keyword": "万户OA"
		}]
	},
	{
		"cms": "悟空CRM",
		"rules": [{
			"type": "title",
			"keyword": "悟空CRM"
		}]
	},
	{
		"cms": "小米路由器",
		"rules": [{
			"type": "title",
			"keyword": "小米路由器\" "
		}]
	},
	{
		"cms": "沃科网异网同显系统",
		"rules": [{
			"type": "title",
			"keyword": "异网同显系统"
		}]
	},
	{
		"cms": "易分析",
		"rules": [{
			"type": "title",
			"keyword": "易分析 PHPStat Analytics"
		}]
	},
	{
		"cms": "用友erp-nc",
		"rules": [{
			"type": "title",
			"keyword": "用友新世纪"
		}]
	},
	{
		"cms": "用友致远oa",
		"rules": [{
			"type": "title",
			"keyword": "用友致远OA\""
		}]
	}
]
`)

var keywordServerSlice = []string{"MochiWeb", "webfs", "WebCam2000", "WebSite", "Statistics Server",
	"Virata", "RemotelyAnywhere", "Ipswitch", "IMail", "Ipswitch Web Calendaring",
	"Check Point SVN foundation", "HP", "Allegro", "dhttpd", "Snap Appliance",
	"Zeus", "IP", "PRINT", "DHost", "3ware",
	"Cherokee", "WindWeb", "SimpleServer", "Xitami", "Resin",
	"linuxconf", "TinyWeb", "WebSitePro", "Lucent Security Management Admin Server", "thttpd",
	"kpf", "dwhttpd", "LispWeb", "WDaemon", "Oracle XML DB",
	"Oracle9iAS", "ALT", "Apt", "Sun", "SunONE WebServer",
	"IBM", "Embedded HTTP Server", "Embedded HTTP Server v", "Embedded HTTP Server USR", "PSIWBL",
	"WhatsUp", "Roxen", "SAMBAR", "15:38", "WebLogic",
	"WebLogic Server", "icecast", "HP Web Jetadmin", "Tomcat Web Server", "Apache",
	"publicfile", "Apache Tomcat", "Stronghold", "nginx", "CommuniGatePro",
	"DSS", "QTSS", "Fnord", "MiniServ", "NetWare HTTP Stack",
	"HTTPd", "Lotus", "Domino", "Insight Manager", "Netscape",
	"NCSA", "Oracle HTTP Server Powered by Apache", "cisco", "NetApp", "Boa",
	"Jetty", "MortBay", "WebSphere Application Server", "JRun Web Server", "RomPager",
	"IDSL MailGate", "Anti", "BaseHTTP", "InterMapper", "ZOT",
	"Lasso", "fwlogwatch", "GNUMP3d", "Gatling", "Viavideo",
	"IBM HTTP Server", "SAlive", "Abyss", "LseriesWeb", "AOLserver",
	"Orion", "Agent", "RMC Webserver", "TwistedWeb", "Twisted",
	"Azureus", "Microsoft", "DManager", "IDS", "Indy",
	"FrontPage", "TSM", "ADSM", "EPSON", "Monkey",
	"Monkey Server", "wr", "FSPMS", "RapidLogic", "TUX",
	"2Wire", "2wire Gateway", "Oracle Application Server Containers for J2EE", "Oracle Application Server Containers for J2EE 10g", "Oracle",
	"OracleAS", "Oracle Containers for J2EE", "Askey Software", "WN", "GST",
	"LANCOM", "BrowseAmp", "Thy", "FileMakerPro", "AdSubtract",
	"bozohttpd", "Null httpd", "Dune", "Meredydd Luff", "zawhttpd",
	"NeepHttpd", "SHS", "cpaneld", "cpsrvd", "CERN",
	"Savant", "TiVo Server", "WebTopia", "Microsoft ASP", "GWS",
	"GFE", "Ubicom", "Boche", "libwww", "MACOS",
	"Sinclair ZX", "WWW File Share Pro", "HP Apache", "Internet Firewall", "swcd",
	"LiveStats Reporting Server", "Embedded HTTPD v", "Spyglass", "Tcl", "ListManagerWeb",
	"4D", "Caudium", "JC", "GeoHttpServer", "ATR",
	"PortWise mVPN", "WYM", "WatchGuard Firewall", "NetPort Software", "OmniHTTPd",
	"OmniSecure", "MirandaWeb", "Mirapoint", "Unknown", "LabVIEW",
	"Groove", "Agranat", "WebSnmp Server Httpd", "Plan9", "IceWarp WebSrv",
	"IceWarp", "BBIagent", "Web", "Rapid Logic", "OpenLink",
	"Embperl", "SiteScope", "SilverStream Server", "PicoWebServer", "tivo",
	"Dahlia", "WEBrick", "Mathopd", "ml", "fhttpd",
	"MyServer", "Quantum Corporation", "Simple java", "Micro", "Mono",
	"SimpleHTTP", "Cougar", "Footprint", "LogMeIn Web Gateway", "ArGoSoft Mail Server Freeware",
	"Fastream NETFile Web Server", "WebServer", "HDS Hi", "WebTrends HTTP Server", "Desktop On",
	"OCServer", "ENI", "JavaWebServer", "whostmgr", "BBC",
	"Servertec", "DirectUpdate", "CCS", "VisiBroker", "Compaq Insight Manager XE",
	"ISS", "Jigsaw", "gnump3d2", "SpyBot", "WWW",
	"Ares", "Paws", "monit", "Red Carpet Daemon", "CL",
	"SAP", "JTALKServer", "iSoft Commerce Suite Server", "MS", "BSE",
	"WebMail", "PRTG", "SmarterTools", "Project Engine Server", "McAfee",
	"HoneydHTTP", "MultiSync Plugin", "C4D", "Sun Java System Application Server Platform Edition", "Sun Java System Application Server",
	"AltaVista Avhttpd", "Servage", "Nucleus WebServ", "Thunderstone", "NETID",
	"MailEnable", "Kerio Embedded WebServer", "and", "Kerio MailServer", "KHAPI",
	"Mediasurface", "SQ", "BeOS", "Grandstream", "LiteServe",
	"YAZ", "svea", "BRS", "eSoft", "Tandberg Television Web server",
	"Novell", "Intoto", "httrack", "GeneWeb", "jabberd",
	"VB150", "iGuard Embedded Web Server", "NetworkActiv", "ATEN HTTP Server", "CPWS",
	"Niagara Web Server", "The Knopflerfish HTTP Server", "HTTP", "Nanox WebServer", "ipMonitor",
	"Tarantella", "Web Crossing", "Kannel", "ZyXEL", "DIONIS",
	"VOMwebserver v", "AKCP Embedded Web Server", "SWS", "Sensorsoft", "Meridian Data",
	"IPCheck", "glass", "Spinnaker", "GoAhead", "SAP Web Application Server",
	"SentinelKeysServer", "Techno Vision Security System Ver", "webcamXP", "Apple Embedded Web Server", "iPrism",
	"XOS", "Systinet Server for Java", "tracd", "Sametime Server", "OCaml HTTP Daemon",
	"Anapod Manager", "IISGuard", "DesktopAuthority", "LogMeIn", "MacroMaker",
	"Siemens Gigaset C450 IP", "nhttpd", "aidex", "Polycom SoundPoint IP Telephone HTTPd", "David",
	"PeopleSoft RENSRV", "HFS", "Ultraseek", "Web Server", "Henry",
	"Freechal P2P", "Httpinfo olsrd plugin", "DirectAdmin Daemon v", "WebSEAL", "VCS",
	"DVSS", "pcastd", "BigFixHTTPServer", "Quick", "SentinelProtectionServer",
	"PowerSchool", "Intrinsyc deviceWEB v", "Hitachi Web Server", "FreeBrowser", "GG",
	"Easy File Sharing Web Server v", "darkstat", "Rumpus", "Medusa", "Sphere V",
	"BlueDragon Server", "vdradmind", "NetEVI", "Java", "LiteSpeed",
	"Tornado", "Web Transaction Server For ClearPath MCP", "AnomicHTTPD", "SnapStream", "Yaws",
	"Centile Embedded HTTPSd server", "PWLib", "ICONAG web server", "KTorrent", "Wildcat",
	"Mongrel", "ChatSpace", "IntellipoolHTTPD", "MX4J", "Vistabox",
	"alevtd", "CosminexusComponentContainer", "JAGeX", "BarracudaHTTP", "SAP Internet Graphics Server",
	"SAP Message Server", "SCO I2O Dialogue Daemon", "Cpanel", "NessusWWW", "TIB",
	"Snug", "Grandstream GXP2000", "D", "FM Web Publishing", "Snakelets",
	"GroupWise MTA", "GroupWise POA", "GroupWise GWIA", "Messenger", "Hunchentoot",
	"AllegroServe", "Hop", "Minix httpd", "Eye", "PWS",
	"BCReport", "FTGate", "Comanche", "MediabolicMWEB", "UltiDev Cassini",
	"Swazoo", "RAID HTTPServer", "Axigen", "lighttpd", "Nano HTTPD library",
	"Enigma2 WebInterface Server", "Winstone Servlet Engine v", "LuCId", "GlassFish", "GlassFish Server Open Source Edition",
	"Sun GlassFish Enterprise Server v", "Sun GlassFish Communications Server", "llink", "OctoWebSvr", "CompaqHTTPServer",
	"BAIDA", "IWeb", "Serv", "ASSP", "ZNC ZNC",
	"ZNC", "Kerio Connect", "IdeaWebServer", "TreeNeWS", "CANON HTTP Server",
	"AVGADMINSERVER", "AVGADMINSERVER64", "Texis", "Linux", "Httpd v",
	"Firefly", "HTTP Server", "Zervit", "NGAMS", "Embedthis",
	"EWS", "CubeCoders", "cc", "klhttpd", "Monitorix HTTP Server",
	"YTS", "alphapd", "NG", "Avaya Push Agent Ver x", "NSC",
	"Cosminexus HTTP Server", "Intel", "DrWebServer", "WebIOPi", "Hikvision",
	"switch", "akstreamer", "Netgem", "User Agent Web Server", "tr069 http server",
	"WSTL CPE 1", "DVRDVS", "Ericom Access Server x64", "Ericom Access Server", "xxxxxxxx",
	"RealTimes Desktop Service", "EgdLws", "darkhttpd", "ESTOS WebServer", "Axence nVision WebAccess HTTP Server",
	"ADB Broadband HTTP Server", "ZTE Web Server", "Network", "lwIP", "ghs",
	"ATCOM", "WatchGuard", "calibre", "Mbedthis", "Tntnet",
	"PasteWSGIServer", "FlashCom", "thin", "Conexant", "CherryPy",
	"NetBox Version", "OmikronHTTPOrigin", "Zope", "Optenet Web Server", "uClinux",
	"uc", "Perl Dancer", "Perl Dancer2", "Hiawatha v", "TornadoServer",
	"Pegasus", "Werkzeug", "Webduino", "Restlet", "node",
	"corehttp", "ECS", "cloudflare", "GateOne", "Warp",
	"Vorlon SR", "Rocket", "Debian Apt", "mini", "Splunkd",
	"BarracudaServer", "kolibri", "Karrigell", "WSGIServer", "ExtremeWare",
	"ngx", "openresty", "IntelliJ IDEA", "Cowboy", "Xavante",
	"BBVS", "CE", "Play", "IBM Mobile Connect", "Wave World Wide Web Server",
	"MQX HTTPSRV", "MQX HTTP", "Devline Linia Server", "esp8266", "Mojolicious",
	"Caddy", "KS", "Content Gateway Manager", "bfe", "360wzws",
	"instart", "Tengine", "0W", "NetData Embedded HTTP Server", "Digiweb",
	"Wakanda", "Clearswift", "KFWebServer", "Huawei", "Seattle Lab HTTP Server",
	"WindRiver", "Cassini", "ZK Web Server", "WildFly", "Icinga",
	"Motion", "Simple", "Vidat V7", "PowerStudio v", "servX",
	"JREntServer", "Prime", "WebSTAR"}

var keywordXPoweredBySlice = []string{"PHP", "ASP.NET", "ThinkPHP", "topsec"}
