package httpfinger

import "strings"

type keywordFinger []struct {
	Type    string `json:"type"`
	Cms     string `json:"cms"`
	Keyword string `json:"keyword"`
}

var KeywordFinger keywordFinger

func (k keywordFinger) Match(body string, header string) string {
	for _, kSub := range k {
		if kSub.Type == "body" {
			if strings.Contains(body, kSub.Keyword) {
				return kSub.Cms
			}
		}
		if kSub.Type == "header" {
			if strings.Contains(header, kSub.Keyword) {
				return kSub.Cms
			}
		}
	}
	return ""
}

var keywordFingerByte = []byte(`
[{
		"type": "body",
		"cms": "Dell-Printer-",
		"keyword": "title=\"Dell Laser Printer\""
	},
	{
		"type": "body",
		"cms": "HP-OfficeJet-Printer-",
		"keyword": "title=\"HP Officejet\" || body=\"align=\"center\">HP Officejet\""
	},
	{
		"type": "body",
		"cms": "Biscom-Delivery-Server-",
		"keyword": "body=\"/bds/stylesheets/fds.css\" || body=\"/bds/includes/fdsJavascript.do\""
	},
	{
		"type": "body",
		"cms": "DD-WRT-",
		"keyword": "body=\"style/pwc/ddwrt.css\""
	},
	{
		"type": "body",
		"cms": "ewebeditor-",
		"keyword": "body=\"/ewebeditor.htm?\""
	},
	{
		"type": "body",
		"cms": "fckeditor-",
		"keyword": "body=\"new FCKeditor\""
	},
	{
		"type": "body",
		"cms": "xheditor-",
		"keyword": "body=\"xheditor_lang/zh-cn.js\"||body=\"class=\"xheditor\"||body=\".xheditor(\""
	},
	{
		"type": "body",
		"cms": "百为路由-",
		"keyword": "body=\"提交验证的id必须是ctl_submit\""
	},
	{
		"type": "body",
		"cms": "锐捷NBR路由器-",
		"keyword": "body=\"free_nbr_login_form.png\""
	},
	{
		"type": "body",
		"cms": "mikrotik-",
		"keyword": "title=\"RouterOS\" && body=\"mikrotik\""
	},
	{
		"type": "body",
		"cms": "ThinkSNS-",
		"keyword": "body=\"/addons/theme/\" && body=\"全局变量\""
	},
	{
		"type": "body",
		"cms": "h3c路由器-",
		"keyword": "title=\"Web user login\" && body=\"nLanguageSupported\""
	},
	{
		"type": "body",
		"cms": "jcg无线路由器-",
		"keyword": "title=\"Wireless Router\" && body=\"http://www.jcgcn.com\""
	},
	{
		"type": "body",
		"cms": "D-Link_VoIP_Wireless_Router-",
		"keyword": "title=\"D-Link VoIP Wireless Router\""
	},
	{
		"type": "body",
		"cms": "arrisi_Touchstone-",
		"keyword": "title=\"Touchstone Status\" || body=\"passWithWarnings\""
	},
	{
		"type": "body",
		"cms": "ZyXEL-",
		"keyword": "body=\"Forms/rpAuth_1\""
	},
	{
		"type": "body",
		"cms": "Ruckus-",
		"keyword": "body=\"mon.  Tell me your username\" || title=\"Ruckus Wireless Admin\""
	},
	{
		"type": "body",
		"cms": "Motorola_SBG900-",
		"keyword": "title=\"Motorola SBG900\""
	},
	{
		"type": "body",
		"cms": "Wimax_CPE-",
		"keyword": "title=\"Wimax CPE Configuration\""
	},
	{
		"type": "body",
		"cms": "Cisco_Cable_Modem-",
		"keyword": "title=\"Cisco Cable Modem\""
	},
	{
		"type": "body",
		"cms": "Scientific-Atlanta_Cable_Modem-",
		"keyword": "title=\"Scientific-Atlanta Cable Modem\""
	},
	{
		"type": "body",
		"cms": "rap-",
		"keyword": "body=\"/jscripts/rap_util.js\""
	},
	{
		"type": "body",
		"cms": "ZTE_MiFi_UNE-",
		"keyword": "title=\"MiFi UNE 4G LTE\""
	},
	{
		"type": "body",
		"cms": "ZTE_ZSRV2_Router-",
		"keyword": "title=\"ZSRV2路由器Web管理系统\" && body=\"ZTE Corporation. All Rights Reserved.\""
	},
	{
		"type": "body",
		"cms": "百为智能流控路由器-",
		"keyword": "title=\"BYTEVALUE 智能流控路由器\" && body=\"<a href=\"http://www.bytevalue.com/\" target=\"_blank\">\""
	},
	{
		"type": "body",
		"cms": "乐视路由器-",
		"keyword": "title=\"乐视路由器\" && body=\"<div class=\"login-logo\"></div>\""
	},
	{
		"type": "body",
		"cms": "Verizon_Wireless_Router-",
		"keyword": "title=\"Wireless Broadband Router Management Console\" && body = \"verizon_logo_blk.gif\""
	},
	{
		"type": "body",
		"cms": "Nexus_NX_router-",
		"keyword": "body=\"http://nexuswifi.com/\" && title=\"Nexus NX\""
	},
	{
		"type": "body",
		"cms": "Verizon_Router-",
		"keyword": "title=\"Verizon Router\""
	},
	{
		"type": "body",
		"cms": "小米路由器-",
		"keyword": "title=\"小米路由器\" "
	},
	{
		"type": "body",
		"cms": "QNO_Router-",
		"keyword": "body=\"/QNOVirtual_Keyboard.js\" && body=\"/images/login_img01_03.gif\""
	},
	{
		"type": "body",
		"cms": "爱快流控路由-",
		"keyword": "title=\"爱快\" && body=\"/resources/images/land_prompt_ico01.gif\""
	},
	{
		"type": "body",
		"cms": "Django-",
		"keyword": "body=\"__admin_media_prefix__\" || body=\"csrfmiddlewaretoken\""
	},
	{
		"type": "body",
		"cms": "axis2-web-",
		"keyword": "body=\"axis2-web/css/axis-style.css\""
	},
	{
		"type": "body",
		"cms": "Apache-Wicket-",
		"keyword": "body=\"xmlns:wicket=\" || body=\"/org.apache.wicket.\""
	},
	{
		"type": "body",
		"cms": "BEA-WebLogic-Server-",
		"keyword": "body=\"<h1>BEA WebLogic Server\" || body=\"WebLogic\""
	},
	{
		"type": "body",
		"cms": "EDK-",
		"keyword": "body=\"<!-- /killlistable.tpl -->\""
	},
	{
		"type": "body",
		"cms": "eDirectory-",
		"keyword": "body=\"target=\"_blank\">eDirectory&trade\" || body=\"Powered by <a href=\"http://www.edirectory.com\""
	},
	{
		"type": "body",
		"cms": "Esvon-Classifieds-",
		"keyword": "body=\"Powered by Esvon\""
	},
	{
		"type": "body",
		"cms": "Fluid-Dynamics-Search-Engine-",
		"keyword": "body=\"content=\"fluid dynamics\""
	},
	{
		"type": "body",
		"cms": "mongodb-",
		"keyword": "body=\"<a href=\"/_replSet\">Replica set status</a></p>\""
	},
	{
		"type": "body",
		"cms": "MVB2000-",
		"keyword": "title=\"MVB2000\" || body=\"The Magic Voice Box\""
	},
	{
		"type": "body",
		"cms": "GPSweb-",
		"keyword": "title=\"GPSweb\""
	},
	{
		"type": "body",
		"cms": "phpinfo-",
		"keyword": "title=\"phpinfo\" && body=\"Virtual Directory Support \""
	},
	{
		"type": "body",
		"cms": "lemis管理系统-",
		"keyword": "body=\"lemis.WEB_APP_NAME\""
	},
	{
		"type": "body",
		"cms": "FreeboxOS-",
		"keyword": "title=\"Freebox OS\" || body=\"logo_freeboxos\""
	},
	{
		"type": "body",
		"cms": "Wimax_CPE-",
		"keyword": "title=\"Wimax CPE Configuration\""
	},
	{
		"type": "body",
		"cms": "Scientific-Atlanta_Cable_Modem-",
		"keyword": "title=\"Scientific-Atlanta Cable Modem\""
	},
	{
		"type": "body",
		"cms": "rap-",
		"keyword": "body=\"/jscripts/rap_util.js\""
	},
	{
		"type": "body",
		"cms": "ZTE_MiFi_UNE-",
		"keyword": "title=\"MiFi UNE 4G LTE\""
	},
	{
		"type": "body",
		"cms": "用友商战实践平台-",
		"keyword": "body=\"Login_Main_BG\" && body=\"Login_Owner\""
	},
	{
		"type": "body",
		"cms": "moosefs-",
		"keyword": "body=\"mfs.cgi\" || body=\"under-goal files\""
	},
	{
		"type": "body",
		"cms": "蓝盾BDWebGuard-",
		"keyword": "body=\"BACKGROUND: url(images/loginbg.jpg) #e5f1fc\""
	},
	{
		"type": "body",
		"cms": "护卫神网站安全系统-",
		"keyword": "title=\"护卫神.网站安全系统\""
	},
	{
		"type": "body",
		"cms": "phpDocumentor-",
		"keyword": "body=\"Generated by phpDocumentor\""
	},
	{
		"type": "body",
		"cms": "Adobe_ CQ5-",
		"keyword": "body=\"_jcr_content\""
	},
	{
		"type": "body",
		"cms": "Adobe_GoLive-",
		"keyword": "body=\"generator\" content=\"Adobe GoLive\""
	},
	{
		"type": "body",
		"cms": "Adobe_RoboHelp-",
		"keyword": "body=\"generator\" content=\"Adobe RoboHelp\""
	},
	{
		"type": "body",
		"cms": "Amaya-",
		"keyword": "body=\"generator\" content=\"Amaya\""
	},
	{
		"type": "body",
		"cms": "OpenMas-",
		"keyword": "title=\"OpenMas\" || body=\"loginHead\"><link href=\"App_Themes\""
	},
	{
		"type": "body",
		"cms": "recaptcha-",
		"keyword": "body=\"recaptcha_ajax.js\""
	},
	{
		"type": "body",
		"cms": "TerraMaster-",
		"keyword": "title=\"TerraMaster\" && body=\"/js/common.js\""
	},
	{
		"type": "body",
		"cms": "创星伟业校园网群-",
		"keyword": "body=\"javascripts/float.js\" && body=\"vcxvcxv\""
	},
	{
		"type": "body",
		"cms": "正方教务管理系统-",
		"keyword": "body=\"style/base/jw.css\""
	},
	{
		"type": "body",
		"cms": "UFIDA_NC-",
		"keyword": "(body=\"UFIDA\" && body=\"logo/images/\") || body=\"logo/images/ufida_nc.png\""
	},
	{
		"type": "body",
		"cms": "北创图书检索系统-",
		"keyword": "body=\"opac_two\""
	},
	{
		"type": "body",
		"cms": "北京清科锐华CEMIS-",
		"keyword": "body=\"/theme/2009/image\" && body=\"login.asp\""
	},
	{
		"type": "body",
		"cms": "RG-PowerCache内容加速系统-",
		"keyword": "title=\"RG-PowerCache\""
	},
	{
		"type": "body",
		"cms": "sugon_gridview-",
		"keyword": "body=\"/common/resources/images/common/app/gridview.ico\""
	},
	{
		"type": "body",
		"cms": "SLTM32_Configuration-",
		"keyword": "title=\"SLTM32 Web Configuration Pages \""
	},
	{
		"type": "body",
		"cms": "SHOUTcast-",
		"keyword": "title=\"SHOUTcast Administrator\""
	},
	{
		"type": "body",
		"cms": "milu_seotool-",
		"keyword": "body=\"plugin.php?id=milu_seotool\""
	},
	{
		"type": "body",
		"cms": "CISCO_EPC3925-",
		"keyword": "body=\"Docsis_system\" && body=\"EPC3925\""
	},
	{
		"type": "body",
		"cms": "HP_iLO(HP_Integrated_Lights-Out)-",
		"keyword": "body=\"js/iLO.js\""
	},
	{
		"type": "body",
		"cms": "Siemens_SIMATIC-",
		"keyword": "body=\"/S7Web.css\""
	},
	{
		"type": "body",
		"cms": "Schneider_Quantum_140NOE77101-",
		"keyword": "body=\"indexLanguage\" && body=\"html/config.js\""
	},
	{
		"type": "body",
		"cms": "lynxspring_JENEsys-",
		"keyword": "body=\"LX JENEsys\""
	},
	{
		"type": "body",
		"cms": "Sophos_Web_Appliance-",
		"keyword": "title=\"Sophos Web Appliance\""
	},
	{
		"type": "body",
		"cms": "Comcast_Business-",
		"keyword": "body=\"cmn/css/common-min.css\""
	},
	{
		"type": "body",
		"cms": "Locus_SolarNOC-",
		"keyword": "title=\"SolarNOC - Login\""
	},
	{
		"type": "body",
		"cms": "Everything-",
		"keyword": "(body=\"Everything.gif\"||body=\"everything.png\") && title=\"Everything\""
	},
	{
		"type": "body",
		"cms": "honeywell NetAXS-",
		"keyword": "title=\"Honeywell NetAXS\""
	},
	{
		"type": "body",
		"cms": "Symantec Messaging Gateway-",
		"keyword": "title=\"Messaging Gateway\""
	},
	{
		"type": "body",
		"cms": "xfinity-",
		"keyword": "title=\"Xfinity\" || body=\"/reset-meyer-1.0.min.css\""
	},
	{
		"type": "body",
		"cms": "网动云视讯平台-",
		"keyword": "title=\"Acenter\" || body=\"/js/roomHeight.js\" || body=\"meetingShow!show.action\""
	},
	{
		"type": "body",
		"cms": "蓝凌EIS智慧协同平台-",
		"keyword": "body=\"/scripts/jquery.landray.common.js\" || body=\"v11_QRcodeBar clr\""
	},
	{
		"type": "body",
		"cms": "金山KingGate-",
		"keyword": "body=\"/src/system/login.php\""
	},
	{
		"type": "body",
		"cms": "天融信入侵检测系统TopSentry-",
		"keyword": "title=\"天融信入侵检测系统TopSentry\""
	},
	{
		"type": "body",
		"cms": "天融信日志收集与分析系统-",
		"keyword": "title=\"天融信日志收集与分析系统\""
	},
	{
		"type": "body",
		"cms": "天融信WEB应用防火墙-",
		"keyword": "title=\"天融信WEB应用防火墙\""
	},
	{
		"type": "body",
		"cms": "天融信入侵防御系统TopIDP-",
		"keyword": "body=\"天融信入侵防御系统TopIDP\""
	},
	{
		"type": "body",
		"cms": "天融信Web应用安全防护系统-",
		"keyword": "title=\"天融信Web应用安全防护系统\""
	},
	{
		"type": "body",
		"cms": "天融信TopFlow-",
		"keyword": "body=\"天融信TopFlow\""
	},
	{
		"type": "body",
		"cms": "汉码软件-",
		"keyword": "title=\"汉码软件\" || body=\"alt=\"汉码软件LOGO\" || body=\"content=\"汉码软件\""
	},
	{
		"type": "body",
		"cms": "凡科-",
		"keyword": "body=\"凡科互联网科技股份有限公司\" || body=\"content=\"凡科\""
	},
	{
		"type": "body",
		"cms": "易分析-",
		"keyword": "title=\"易分析 PHPStat Analytics\" || body=\"PHPStat Analytics 网站数据分析系统\""
	},
	{
		"type": "body",
		"cms": "phpems考试系统-",
		"keyword": "title=\"phpems\" || body=\"content=\"PHPEMS\""
	},
	{
		"type": "body",
		"cms": "智睿软件-",
		"keyword": "body=\"content=\"智睿软件\" || body=\"Zhirui.js\""
	},
	{
		"type": "body",
		"cms": "Apabi数字资源平台-",
		"keyword": "body=\"Default/apabi.css\" || body=\"<link href=\"HTTP://apabi\" || title=\"数字资源平台\""
	},
	{
		"type": "body",
		"cms": "Fortinet Firewall-",
		"keyword": "title=\"Firewall Notification\""
	},
	{
		"type": "body",
		"cms": "WDlinux-",
		"keyword": "title=\"wdOS\""
	},
	{
		"type": "body",
		"cms": "小脑袋-",
		"keyword": "body=\"http://stat.xiaonaodai.com/stat.php\""
	},
	{
		"type": "body",
		"cms": "天融信ADS管理平台-",
		"keyword": "title=\"天融信ADS管理平台\""
	},
	{
		"type": "body",
		"cms": "天融信异常流量管理与抗拒绝服务系统-",
		"keyword": "title=\"天融信异常流量管理与抗拒绝服务系统\""
	},
	{
		"type": "body",
		"cms": "天融信网络审计系统-",
		"keyword": "body=\"onclick=\"dlg_download()\""
	},
	{
		"type": "body",
		"cms": "天融信脆弱性扫描与管理系统-",
		"keyword": "title=\"天融信脆弱性扫描与管理系统\" || body=\"/js/report/horizontalReportPanel.js\""
	},
	{
		"type": "body",
		"cms": "AllNewsManager_NET-",
		"keyword": "body=\"Powered by\" && body=\"AllNewsManager\""
	},
	{
		"type": "body",
		"cms": "Advanced-Image-Hosting-Script-",
		"keyword": "(body=\"Powered by\" && body=\"yabsoft.com\" ) || body=\"Welcome to install AIHS Script\""
	},
	{
		"type": "body",
		"cms": "SNB股票交易软件-",
		"keyword": "body=\"Copyright 2005–2009 <a href=\"http://www.s-mo.com\">\""
	},
	{
		"type": "body",
		"cms": "AChecker Web accessibility evaluation tool-",
		"keyword": "body=\"content=\"AChecker is a Web accessibility\" || title=\"Checker : Web Accessibility Checker\""
	},
	{
		"type": "body",
		"cms": "SCADA PLC-",
		"keyword": "body=\"/images/rockcolor.gif\" || body=\"/ralogo.gif\" || body=\"Ethernet Processor\""
	},
	{
		"type": "body",
		"cms": ".NET-",
		"keyword": "body=\"content=\"Visual Basic .NET 7.1\""
	},
	{
		"type": "body",
		"cms": "phpmoadmin-",
		"keyword": "title=\"phpmoadmin\""
	},
	{
		"type": "body",
		"cms": "SOMOIDEA-",
		"keyword": "body=\"DESIGN BY SOMOIDEA\""
	},
	{
		"type": "body",
		"cms": "Apache-Archiva-",
		"keyword": "title=\"Apache Archiva\" || body=\"/archiva.js\" || body=\"/archiva.css\""
	},
	{
		"type": "body",
		"cms": "AM4SS-",
		"keyword": "body=\"Powered by am4ss\" || body=\"am4ss.css\""
	},
	{
		"type": "body",
		"cms": "ASPThai_Net-Webboard-",
		"keyword": "body=\"ASPThai.Net Webboard\""
	},
	{
		"type": "body",
		"cms": "Astaro-Command-Center-",
		"keyword": "body=\"/acc_aggregated_reporting.js\" || body=\"/js/_variables_from_backend.js?\""
	},
	{
		"type": "body",
		"cms": "ASP-Nuke-",
		"keyword": "body=\"CONTENT=\"ASP-Nuke\" || body=\"content=\"ASPNUKE\""
	},
	{
		"type": "body",
		"cms": "ASProxy-",
		"keyword": "body=\"Surf the web invisibly using ASProxy power\" || body=\"btnASProxyDisplayButton\""
	},
	{
		"type": "body",
		"cms": "ashnews-",
		"keyword": "body=\"powered by\" && body=\"ashnews\""
	},
	{
		"type": "body",
		"cms": "Arab-Portal-",
		"keyword": "body=\"Powered by: Arab\""
	},
	{
		"type": "body",
		"cms": "AppServ-",
		"keyword": "body=\"appserv/softicon.gif\" || body=\"index.php?appservlang=th\""
	},
	{
		"type": "body",
		"cms": "VZPP Plesk-",
		"keyword": "title=\"VZPP Plesk \""
	},
	{
		"type": "body",
		"cms": "ApPHP-Calendar-",
		"keyword": "body=\"This script was generated by ApPHP Calendar\""
	},
	{
		"type": "body",
		"cms": "BigDump-",
		"keyword": "title=\"BigDump\" || body=\"BigDump: Staggered MySQL Dump Importer\""
	},
	{
		"type": "body",
		"cms": "BestShopPro-",
		"keyword": "body=\"content=\"www.bst.pl\""
	},
	{
		"type": "body",
		"cms": "BASE-",
		"keyword": "body=\"<!-- Basic Analysis and Security Engine (BASE) -->\" || body=\"mailto:base@secureideas.net\""
	},
	{
		"type": "body",
		"cms": "Basilic-",
		"keyword": "body=\"/Software/Basilic\""
	},
	{
		"type": "body",
		"cms": "Basic-PHP-Events-Lister-",
		"keyword": "body=\"Powered by: <a href=\"http://www.mevin.com/\">\""
	},
	{
		"type": "body",
		"cms": "AV-Arcade-",
		"keyword": "body=\"Powered by <a href=\"http://www.avscripts.net/avarcade/\""
	},
	{
		"type": "body",
		"cms": "Auxilium-PetRatePro-",
		"keyword": "body=\"index.php?cmd=11\""
	},
	{
		"type": "body",
		"cms": "Atomic-Photo-Album-",
		"keyword": "body=\"Powered by\" && body=\"Atomic Photo Album\""
	},
	{
		"type": "body",
		"cms": "Axis-PrintServer-",
		"keyword": "body=\"psb_printjobs.gif\" || body=\"/cgi-bin/prodhelp?prod=\""
	},
	{
		"type": "body",
		"cms": "TeamViewer-",
		"keyword": "body=\"This site is running\"&&body=\"TeamViewer\""
	},
	{
		"type": "body",
		"cms": "BlueQuartz-",
		"keyword": "body=\"VALUE=\"Copyright (C) 2000, Cobalt Networks\" || title=\"Login - BlueQuartz\""
	},
	{
		"type": "body",
		"cms": "BlueOnyx-",
		"keyword": "title=\"Login - BlueOnyx\" || body=\"Thank you for using the BlueOnyx\""
	},
	{
		"type": "body",
		"cms": "BMC-Remedy-",
		"keyword": "title=\"Remedy Mid Tier\""
	},
	{
		"type": "body",
		"cms": "BM-Classifieds-",
		"keyword": "body=\"<!-- START HEADER TABLE - HOLDS GRAPHIC AND SITE NAME -->\""
	},
	{
		"type": "body",
		"cms": "Citrix-Metaframe-",
		"keyword": "body=\"window.location=\"/Citrix/MetaFrame\""
	},
	{
		"type": "body",
		"cms": "Cogent-DataHub-",
		"keyword": "body=\"/images/Cogent.gif\" || title=\"Cogent DataHub WebView\""
	},
	{
		"type": "body",
		"cms": "ClipShare-",
		"keyword": "body=\"<!--!!!!!!!!!!!!!!!!!!!!!!!!! Processing SCRIPT\" || body=\"Powered By <a href=\"http://www.clip-share.com\""
	},
	{
		"type": "body",
		"cms": "CGIProxy-",
		"keyword": "body=\"<a href=\"http://www.jmarshall.com/tools/cgiproxy/\""
	},
	{
		"type": "body",
		"cms": "CF-Image-Hosting-Script-",
		"keyword": "body=\"Powered By <a href=\"http://codefuture.co.uk/projects/imagehost/\""
	},
	{
		"type": "body",
		"cms": "Censura-",
		"keyword": "body=\"Powered by: <a href=\"http://www.censura.info\""
	},
	{
		"type": "body",
		"cms": "CA-SiteMinder-",
		"keyword": "body=\"<!-- SiteMinder Encoding\""
	},
	{
		"type": "body",
		"cms": "Carrier-CCNWeb-",
		"keyword": "body=\"/images/CCNWeb.gif\" || body=\"<APPLET CODE=\"JLogin.class\" ARCHIVE=\"JLogin.jar\""
	},
	{
		"type": "body",
		"cms": "cInvoice-",
		"keyword": "body=\"Powered by <a href=\"http://www.forperfect.com/\""
	},
	{
		"type": "body",
		"cms": "Bomgar-",
		"keyword": "body=\"alt=\"Remote Support by BOMGAR\" || body=\"<a href=\"http://www.bomgar.com/products\" class=\"inverse\""
	},
	{
		"type": "body",
		"cms": "cApexWEB-",
		"keyword": "body=\"/capexweb.parentvalidatepassword\" || body=\"name=\"dfparentdb\""
	},
	{
		"type": "body",
		"cms": "CameraLife-",
		"keyword": "body=\"content=\"Camera Life\" || body=\"This site is powered by Camera Life\""
	},
	{
		"type": "body",
		"cms": "CalendarScript-",
		"keyword": "title=\"Calendar Administration : Login\" || body=\"Powered by <A HREF=\"http://www.CalendarScript.com\""
	},
	{
		"type": "body",
		"cms": "Cachelogic-Expired-Domains-Script-",
		"keyword": "body=\"href=\"http://cachelogic.net\">Cachelogic.net\""
	},
	{
		"type": "body",
		"cms": "Burning-Board-Lite-",
		"keyword": "body=\"Powered by <b><a href=\"http://www.woltlab.de\" || body=\"Powered by <b>Burning Board\""
	},
	{
		"type": "body",
		"cms": "Buddy-Zone-",
		"keyword": "body=\"Powered By <a href=\"http://www.vastal.com\" || body=\">Buddy Zone</a>\""
	},
	{
		"type": "body",
		"cms": "Bulletlink-Newspaper-Template-",
		"keyword": "body=\"/ModalPopup/core-modalpopup.css\" || body=\"powered by bulletlink\""
	},
	{
		"type": "body",
		"cms": "Brother-Printer-",
		"keyword": "body=\"<FRAME SRC=\"/printer/inc_head.html\" || body=\"<IMG src=\"/common/image/HL4040CN\""
	},
	{
		"type": "body",
		"cms": "Daffodil-CRM-",
		"keyword": "body=\"Powered by Daffodil\" || body=\"Design & Development by Daffodil Software Ltd\""
	},
	{
		"type": "body",
		"cms": "Cyn_in-",
		"keyword": "body=\"content=\"cyn.in\" || body=\"Powered by cyn.in\""
	},
	{
		"type": "body",
		"cms": "Oracle_OPERA-",
		"keyword": "title=\"MICROS Systems Inc., OPERA\" || body=\"OperaLogin/Welcome.do\""
	},
	{
		"type": "body",
		"cms": "DUgallery-",
		"keyword": "body=\"Powered by DUportal\" ||  body=\"DUgallery\""
	},
	{
		"type": "body",
		"cms": "DublinCore-",
		"keyword": "body=\"name=\"DC.title\""
	},
	{
		"type": "body",
		"cms": "DZCP-",
		"keyword": "body=\"<!--[ DZCP\""
	},
	{
		"type": "body",
		"cms": "DVWA-",
		"keyword": "title=\"Damn Vulnerable Web App (DVWA) - Login\" || body=\"dvwa/css/login.css\" || body=\"dvwa/images/login_logo.png\""
	},
	{
		"type": "body",
		"cms": "DORG-",
		"keyword": "title=\"DORG - \" || body=\"CONTENT=\"DORG\""
	},
	{
		"type": "body",
		"cms": "VOS3000-",
		"keyword": "title=\"VOS3000\"||body=\"<meta name=\"keywords\" content=\"VOS3000\"||body=\"<meta name=\"description\" content=\"VOS3000\"||body=\"images/vos3000.ico\""
	},
	{
		"type": "body",
		"cms": "Elite-Gaming-Ladders-",
		"keyword": "body=\"Powered by Elite\""
	},
	{
		"type": "body",
		"cms": "Entrans-",
		"keyword": "title=\"Entrans\""
	},
	{
		"type": "body",
		"cms": "GateQuest-PHP-Site-Recommender-",
		"keyword": "title=\"GateQuest\""
	},
	{
		"type": "body",
		"cms": "Gallarific-",
		"keyword": "body=\"content=\"Gallarific\" || title=\"Gallarific > Sign in\""
	},
	{
		"type": "body",
		"cms": "EZCMS-",
		"keyword": "body=\"Powered by EZCMS\" || body=\"EZCMS Content Management System\""
	},
	{
		"type": "body",
		"cms": "Etano-",
		"keyword": "body=\"Powered by <a href=\"http://www.datemill.com\" || body=\"Etano</a>. All Rights Reserved.\""
	},
	{
		"type": "body",
		"cms": "GeoServer-",
		"keyword": "body=\"/org.geoserver.web.GeoServerBasePage/\" || body=\"class=\"geoserver lebeg\""
	},
	{
		"type": "body",
		"cms": "GeoNode-",
		"keyword": "body=\"Powered by <a href=\"http://geonode.org\" || body=\"href=\"/catalogue/opensearch\" title=\"GeoNode Search\""
	},
	{
		"type": "body",
		"cms": "Help-Desk-Software-",
		"keyword": "body=\"target=\"_blank\">freehelpdesk.org\""
	},
	{
		"type": "body",
		"cms": "GridSite-",
		"keyword": "body=\"<a href=\"http://www.gridsite.org/\">GridSite\" || body=\"gridsite-admin.cgi?cmd\""
	},
	{
		"type": "body",
		"cms": "GenOHM-SCADA-",
		"keyword": "title=\"GenOHM Scada Launcher\" || body=\"/cgi-bin/scada-vis/\""
	},
	{
		"type": "body",
		"cms": "Infomaster-",
		"keyword": "body=\"/MasterView.css\" || body=\"/masterView.js\" || body=\"/MasterView/MPLeftNavStyle/PanelBar.MPIFMA.css\""
	},
	{
		"type": "body",
		"cms": "Imageview-",
		"keyword": "body=\"content=\"Imageview\" || body=\"By Jorge Schrauwen\" || body=\"href=\"http://www.blackdot.be\" title=\"Blackdot.be\""
	},
	{
		"type": "body",
		"cms": "Ikonboard-",
		"keyword": "body=\"content=\"Ikonboard\" || body=\"Powered by <a href=\"http://www.ikonboard.com\""
	},
	{
		"type": "body",
		"cms": "i-Gallery-",
		"keyword": "title=\"i-Gallery\" || body=\"href=\"igallery.asp\""
	},
	{
		"type": "body",
		"cms": "OrientDB-",
		"keyword": "title=\"Redirecting to OrientDB\""
	},
	{
		"type": "body",
		"cms": "Solr-",
		"keyword": "title=\"Solr Admin\"||body=\"SolrCore Initialization Failures\"||body=\"app_config.solr_path\""
	},
	{
		"type": "body",
		"cms": "Inout-Adserver-",
		"keyword": "body=\"Powered by Inoutscripts\""
	},
	{
		"type": "body",
		"cms": "ionCube-Loader-",
		"keyword": "body=\"alt=\"ionCube logo\""
	},
	{
		"type": "body",
		"cms": "Jamroom-",
		"keyword": "body=\"content=\"Talldude Networks\" || body=\"content=\"Jamroom\""
	},
	{
		"type": "body",
		"cms": "Juniper-NetScreen-Secure-Access-",
		"keyword": "body=\"/dana-na/auth/welcome.cgi\""
	},
	{
		"type": "body",
		"cms": "Jcow-",
		"keyword": "body=\"content=\"Jcow\" || body=\"content=\"Powered by Jcow\" || body=\"end jcow_application_box\""
	},
	{
		"type": "body",
		"cms": "InvisionPowerBoard-",
		"keyword": "body=\"Powered by <a href=\"http://www.invisionboard.com\""
	},
	{
		"type": "body",
		"cms": "teamportal-",
		"keyword": "body=\"TS_expiredurl\""
	},
	{
		"type": "body",
		"cms": "VisualSVN-",
		"keyword": "title=\"VisualSVN Server\""
	},
	{
		"type": "body",
		"cms": "Redmine-",
		"keyword": "body=\"Redmine\" && body=\"authenticity_token\""
	},
	{
		"type": "body",
		"cms": "testlink-",
		"keyword": "body=\"testlink_library.js\""
	},
	{
		"type": "body",
		"cms": "mantis-",
		"keyword": "body=\"browser_search_plugin.php?type=id\" || body=\"MantisBT Team\""
	},
	{
		"type": "body",
		"cms": "Mercurial-",
		"keyword": "title=\"Mercurial repositories index\""
	},
	{
		"type": "body",
		"cms": "activeCollab-",
		"keyword": "body=\"powered by activeCollab\" || body=\"<p id=\"powered_by\"><a href=\"http://www.activecollab.com/\"\""
	},
	{
		"type": "body",
		"cms": "Collabtive-",
		"keyword": "title=\"Login @ Collabtive\""
	},
	{
		"type": "body",
		"cms": "CGI:IRC-",
		"keyword": "title=\"CGI:IRC Login\" || body=\"<!-- This is part of CGI:IRC\" || body=\"<small id=\"ietest\"><a href=\"http://cgiirc.org/\""
	},
	{
		"type": "body",
		"cms": "DotA-OpenStats-",
		"keyword": "body=\"content=\"dota OpenStats\" || body=\"content=\"openstats.iz.rs\""
	},
	{
		"type": "body",
		"cms": "eLitius-",
		"keyword": "body=\"content=\"eLitius\" || body=\"target=\"_blank\" title=\"Affiliate\""
	},
	{
		"type": "body",
		"cms": "gCards-",
		"keyword": "body=\"<a href=\"http://www.gregphoto.net/gcards/index.php\""
	},
	{
		"type": "body",
		"cms": "GpsGate-Server-",
		"keyword": "title=\"GpsGate Server - \""
	},
	{
		"type": "body",
		"cms": "iScripts-ReserveLogic-",
		"keyword": "body=\"Powered by <a href=\"http://www.iscripts.com/reservelogic/\""
	},
	{
		"type": "body",
		"cms": "jobberBase-",
		"keyword": "body=\"powered by\" && body=\"http://www.jobberbase.com\" || body=\"Jobber.PerformSearch\" || body=\"content=\"Jobberbase\""
	},
	{
		"type": "body",
		"cms": "LuManager-",
		"keyword": "title=\"LuManager\""
	},
	{
		"type": "body",
		"cms": "主机宝-",
		"keyword": "body=\"您访问的是主机宝服务器默认页\""
	},
	{
		"type": "body",
		"cms": "wdcp管理系统-",
		"keyword": "title=\"wdcp服务器\" || title=\"lanmp_wdcp 安装成功\""
	},
	{
		"type": "body",
		"cms": "LANMP一键安装包-",
		"keyword": "title=\"LANMP一键安装包\""
	},
	{
		"type": "body",
		"cms": "UPUPW-",
		"keyword": "title=\"UPUPW环境集成包\""
	},
	{
		"type": "body",
		"cms": "wamp-",
		"keyword": "title=\"WAMPSERVER\""
	},
	{
		"type": "body",
		"cms": "easypanel-",
		"keyword": "body=\"/vhost/view/default/style/login.css\""
	},
	{
		"type": "body",
		"cms": "awstats_admin-",
		"keyword": "body=\"generator\" content=\"AWStats\" || body=\"<frame name=\"mainleft\" src=\"awstats.pl?config=\""
	},
	{
		"type": "body",
		"cms": "awstats-",
		"keyword": "body=\"awstats.pl?config=\""
	},
	{
		"type": "body",
		"cms": "moosefs-",
		"keyword": "body=\"mfs.cgi\" || body=\"under-goal files\""
	},
	{
		"type": "body",
		"cms": "护卫神主机管理-",
		"keyword": "title=\"护卫神·主机管理系统\""
	},
	{
		"type": "body",
		"cms": "bacula-web-",
		"keyword": "title=\"Webacula\" || title=\"Bacula Web\" || title=\"Bacula-Web\" || title=\"bacula-web\""
	},
	{
		"type": "body",
		"cms": "Webmin-",
		"keyword": "title=\"Login to Webmin\" || body=\"Webmin server on\""
	},
	{
		"type": "body",
		"cms": "Synology_DiskStation-",
		"keyword": "title=\"Synology DiskStation\" || body=\"SYNO.SDS.Session\""
	},
	{
		"type": "body",
		"cms": "Puppet_Node_Manager-",
		"keyword": "title=\"Puppet Node Manager\""
	},
	{
		"type": "body",
		"cms": "wdcp-",
		"keyword": "title=\"wdcp服务器\""
	},
	{
		"type": "body",
		"cms": "Citrix-XenServer-",
		"keyword": "body=\"Citrix Systems, Inc. XenServer\" || body=\"<a href=\"XenCenterSetup.exe\">XenCenter installer</a>\""
	},
	{
		"type": "body",
		"cms": "DSpace-",
		"keyword": "body=\"content=\"DSpace\" || body=\"<a href=\"http://www.dspace.org\">DSpace Software\""
	},
	{
		"type": "body",
		"cms": "dwr-",
		"keyword": "body=\"/dwr/engine.js\""
	},
	{
		"type": "body",
		"cms": "eXtplorer-",
		"keyword": "title=\"Login - eXtplorer\""
	},
	{
		"type": "body",
		"cms": "File-Upload-Manager-",
		"keyword": "title=\"File Upload Manager\" || body=\"<IMG SRC=\"/images/header.jpg\" ALT=\"File Upload Manager\">\""
	},
	{
		"type": "body",
		"cms": "FileNice-",
		"keyword": "body=\"content=\"the fantabulous mechanical eviltwin code machine\" || body=\"fileNice/fileNice.js\""
	},
	{
		"type": "body",
		"cms": "Glossword-",
		"keyword": "body=\"content=\"Glossword\""
	},
	{
		"type": "body",
		"cms": "IBM-BladeCenter-",
		"keyword": "body=\"/shared/ibmbch.png\" || body=\"/shared/ibmbcs.png\" || body=\"alt=\"IBM BladeCenter\""
	},
	{
		"type": "body",
		"cms": "iLO-",
		"keyword": "body=\"href=\"http://www.hp.com/go/ilo\" || title=\"HP Integrated Lights-Out\""
	},
	{
		"type": "body",
		"cms": "Isolsoft-Support-Center-",
		"keyword": "body=\"Powered by: Support Center\""
	},
	{
		"type": "body",
		"cms": "ISPConfig-",
		"keyword": "body=\"powered by <a HREF=\"http://www.ispconfig.org\""
	},
	{
		"type": "body",
		"cms": "Kleeja-",
		"keyword": "body=\"Powered by Kleeja\""
	},
	{
		"type": "body",
		"cms": "Kloxo-Single-Server-",
		"keyword": "body=\"src=\"/img/hypervm-logo.gif\" || body=\"/htmllib/js/preop.js\" || title=\"HyperVM\""
	},
	{
		"type": "body",
		"cms": "易瑞授权访问系统-",
		"keyword": "body=\"/authjsp/login.jsp\" || body=\"FE0174BB-F093-42AF-AB20-7EC621D10488\""
	},
	{
		"type": "body",
		"cms": "MVB2000-",
		"keyword": "title=\"MVB2000\" || body=\"The Magic Voice Box\""
	},
	{
		"type": "body",
		"cms": "NetShare_VPN-",
		"keyword": "title=\"NetShare\" && title=\"VPN\""
	},
	{
		"type": "body",
		"cms": "pmway_E4_crm-",
		"keyword": "title=\"E4\" && title=\"CRM\""
	},
	{
		"type": "body",
		"cms": "srun3000计费认证系统-",
		"keyword": "title=\"srun3000\""
	},
	{
		"type": "body",
		"cms": "Dolibarr-",
		"keyword": "body=\"Dolibarr Development Team\""
	},
	{
		"type": "body",
		"cms": "Parallels Plesk Panel-",
		"keyword": "body=\"Parallels IP Holdings GmbH\""
	},
	{
		"type": "body",
		"cms": "EasyTrace(botwave)-",
		"keyword": "title=\"EasyTrace\" && body=\"login_page\""
	},
	{
		"type": "body",
		"cms": "管理易-",
		"keyword": "body=\"管理易\" && body=\"minierp\""
	},
	{
		"type": "body",
		"cms": "亿赛通DLP-",
		"keyword": "body=\"CDGServer3\""
	},
	{
		"type": "body",
		"cms": "huawei_auth_server-",
		"keyword": "body=\"75718C9A-F029-11d1-A1AC-00C04FB6C223\""
	},
	{
		"type": "body",
		"cms": "瑞友天翼_应用虚拟化系统 -",
		"keyword": "title=\"瑞友天翼－应用虚拟化系统\""
	},
	{
		"type": "body",
		"cms": "360企业版-",
		"keyword": "body=\"360EntInst\""
	},
	{
		"type": "body",
		"cms": "用友erp-nc-",
		"keyword": "body=\"/nc/servlet/nc.ui.iufo.login.Index\" || title=\"用友新世纪\""
	},
	{
		"type": "body",
		"cms": "Array_Networks_VPN-",
		"keyword": "body=\"an_util.js\""
	},
	{
		"type": "body",
		"cms": "juniper_vpn-",
		"keyword": "body=\"welcome.cgi?p=logo\""
	},
	{
		"type": "body",
		"cms": "CEMIS-",
		"keyword": "body=\"<div id=\"demo\" style=\"overflow:hidden\" && title=\"综合项目管理系统登录\""
	},
	{
		"type": "body",
		"cms": "zenoss-",
		"keyword": "body=\"/zport/dmd/\""
	},
	{
		"type": "body",
		"cms": "OpenMas-",
		"keyword": "title=\"OpenMas\" || body=\"loginHead\"><link href=\"App_Themes\""
	},
	{
		"type": "body",
		"cms": "Ultra_Electronics-",
		"keyword": "body=\"/preauth/login.cgi\" || body=\"/preauth/style.css\""
	},
	{
		"type": "body",
		"cms": "NOALYSS-",
		"keyword": "title=\"NOALYSS\""
	},
	{
		"type": "body",
		"cms": "ALCASAR-",
		"keyword": "body=\"valoriserDiv5\""
	},
	{
		"type": "body",
		"cms": "orocrm-",
		"keyword": "body=\"/bundles/oroui/\""
	},
	{
		"type": "body",
		"cms": "Adiscon_LogAnalyzer-",
		"keyword": "title=\"Adiscon LogAnalyzer\" || (body=\"Adiscon LogAnalyzer\" && body=\"Adiscon GmbH\")"
	},
	{
		"type": "body",
		"cms": "Munin-",
		"keyword": "body=\"Auto-generated by Munin\" || body=\"munin-month.html\""
	},
	{
		"type": "body",
		"cms": "MRTG-",
		"keyword": "body=\"Command line is easier to read using \"View Page Properties\" of your browser\" || title=\"MRTG Index Page\" || body=\"commandline was: indexmaker\""
	},
	{
		"type": "body",
		"cms": "元年财务软件-",
		"keyword": "body=\"yuannian.css\" || body=\"/image/logo/yuannian.gif\""
	},
	{
		"type": "body",
		"cms": "UFIDA_NC-",
		"keyword": "(body=\"UFIDA\" && body=\"logo/images/\") || body=\"logo/images/ufida_nc.png\""
	},
	{
		"type": "body",
		"cms": "Webmin-",
		"keyword": "title=\"Login to Webmin\" || body=\"Webmin server on\""
	},
	{
		"type": "body",
		"cms": "锐捷应用控制引擎-",
		"keyword": "body=\"window.open(\"/login.do\",\"airWin\" || title=\"锐捷应用控制引擎\""
	},
	{
		"type": "body",
		"cms": "Storm-",
		"keyword": "title=\"Storm UI\" || body=\"stormtimestr\""
	},
	{
		"type": "body",
		"cms": "Centreon-",
		"keyword": "body=\"Generator\" content=\"Centreon - Copyright\" || title=\"Centreon - IT & Network Monitoring\""
	},
	{
		"type": "body",
		"cms": "FortiGuard-",
		"keyword": "body=\"FortiGuard Web Filtering\" || title=\"Web Filter Block Override\" || body=\"/XX/YY/ZZ/CI/MGPGHGPGPFGHCDPFGGOGFGEH\""
	},
	{
		"type": "body",
		"cms": "PineApp-",
		"keyword": "title=\"PineApp WebAccess - Login\" || body=\"/admin/css/images/pineapp.ico\""
	},
	{
		"type": "body",
		"cms": "CDR-Stats-",
		"keyword": "title=\"CDR-Stats | Customer Interface\" || body=\"/static/cdr-stats/js/jquery\""
	},
	{
		"type": "body",
		"cms": "GenieATM-",
		"keyword": "title=\"GenieATM\" || body=\"Copyright© Genie Networks Ltd.\" || body=\"defect 3531\""
	},
	{
		"type": "body",
		"cms": "Spark_Worker-",
		"keyword": "title=\"Spark Worker at\""
	},
	{
		"type": "body",
		"cms": "Spark_Master-",
		"keyword": "title=\"Spark Master at\""
	},
	{
		"type": "body",
		"cms": "Kibana-",
		"keyword": "title=\"Kibana\" || body=\"kbnVersion\""
	},
	{
		"type": "body",
		"cms": "UcSTAR-",
		"keyword": "title=\"UcSTAR 管理控制台\""
	},
	{
		"type": "body",
		"cms": "i@Report-",
		"keyword": "body=\"ESENSOFT_IREPORT_SERVER\" || body=\"com.sanlink.server.Login\" || body=\"ireportclient\" || body=\"css/ireport.css\""
	},
	{
		"type": "body",
		"cms": "帕拉迪统一安全管理和综合审计系统-",
		"keyword": "body=\"module/image/pldsec.css\""
	},
	{
		"type": "body",
		"cms": "openEAP-",
		"keyword": "title=\"openEAP_统一登录门户\""
	},
	{
		"type": "body",
		"cms": "Dorado-",
		"keyword": "title=\"Dorado Login Page\""
	},
	{
		"type": "body",
		"cms": "金龙卡金融化一卡通网站查询子系统-",
		"keyword": "title=\"金龙卡金融化一卡通网站查询子系统\" || body=\"location.href=\"homeLogin.action\""
	},
	{
		"type": "body",
		"cms": "一采通-",
		"keyword": "body=\"/custom/GroupNewsList.aspx?GroupId=\""
	},
	{
		"type": "body",
		"cms": "埃森诺网络服务质量检测系统-",
		"keyword": "title=\"埃森诺网络服务质量检测系统 \""
	},
	{
		"type": "body",
		"cms": "惠尔顿上网行为管理系统-",
		"keyword": "body=\"updateLoginPswd.php\" && body=\"PassroedEle\""
	},
	{
		"type": "body",
		"cms": "ACSNO网络探针-",
		"keyword": "title=\"探针管理与测试系统-登录界面\""
	},
	{
		"type": "body",
		"cms": "绿盟下一代防火墙-",
		"keyword": "title=\"NSFOCUS NF\""
	},
	{
		"type": "body",
		"cms": "用友U8-",
		"keyword": "body=\"getFirstU8Accid\""
	},
	{
		"type": "body",
		"cms": "华为（HUAWEI）安全设备-",
		"keyword": "body=\"sweb-lib/resource/\""
	},
	{
		"type": "body",
		"cms": "网神防火墙-",
		"keyword": "title=\"secgate 3600\" || body=\"css/lsec/login.css\""
	},
	{
		"type": "body",
		"cms": "cisco UCM-",
		"keyword": "body=\"/ccmadmin/\" || title=\"Cisco Unified\""
	},
	{
		"type": "body",
		"cms": "panabit智能网关-",
		"keyword": "title=\"panabit\""
	},
	{
		"type": "body",
		"cms": "久其通用财表系统-",
		"keyword": "body=\"<nobr>北京久其软件股份有限公司\" || body=\"/netrep/intf\" || body=\"/netrep/message2/\""
	},
	{
		"type": "body",
		"cms": "soeasy网站集群系统-",
		"keyword": "body=\"EGSS_User\" || title=\"SoEasy网站集群\""
	},
	{
		"type": "body",
		"cms": "畅捷通-",
		"keyword": "title=\"畅捷通\""
	},
	{
		"type": "body",
		"cms": "科来RAS-",
		"keyword": "title=\"科来网络回溯\" || body=\"科来软件 版权所有\" || body=\"i18ninit.min.js\""
	},
	{
		"type": "body",
		"cms": "科迈RAS系统-",
		"keyword": "title=\"科迈RAS\" || body=\"type=\"application/npRas\" || body=\"远程技术支持请求：<a href=\"http://www.comexe.cn\""
	},
	{
		"type": "body",
		"cms": "单点CRM系统-",
		"keyword": "body=\"URL=general/ERP/LOGIN/\" || body=\"content=\"单点CRM系统\" ||title=\"客户关系管理-CRM\""
	},
	{
		"type": "body",
		"cms": "中国期刊先知网-",
		"keyword": "body=\"本系统由<span class=\"STYLE1\" ><a href=\"http://www.firstknow.cn\" || body=\"<img src=\"images/logoknow.png\"\""
	},
	{
		"type": "body",
		"cms": "loyaa信息自动采编系统-",
		"keyword": "body=\"/Loyaa/common.lib.js\""
	},
	{
		"type": "body",
		"cms": "浪潮政务系统-",
		"keyword": "body=\"OnlineQuery/QueryList.aspx\" || title=\"浪潮政务\" || body=\"LangChao.ECGAP.OutPortal\""
	},
	{
		"type": "body",
		"cms": "悟空CRM-",
		"keyword": "title=\"悟空CRM\" || body=\"/Public/js/5kcrm.js\""
	},
	{
		"type": "body",
		"cms": "用友ufida-",
		"keyword": "body=\"/System/Login/Login.asp?AppID=\""
	},
	{
		"type": "body",
		"cms": "金蝶EAS-",
		"keyword": "body=\"easSessionId\""
	},
	{
		"type": "body",
		"cms": "金蝶政务GSiS-",
		"keyword": "body=\"/kdgs/script/kdgs.js\""
	},
	{
		"type": "body",
		"cms": "网御上网行为管理系统-",
		"keyword": "title=\"Leadsec ACM\""
	},
	{
		"type": "body",
		"cms": "ZKAccess 门禁管理系统-",
		"keyword": "body=\"/logoZKAccess_zh-cn.jpg\""
	},
	{
		"type": "body",
		"cms": "福富安全基线管理-",
		"keyword": "body=\"align=\"center\">福富软件\""
	},
	{
		"type": "body",
		"cms": "中控智慧时间安全管理平台-",
		"keyword": "title=\"ZKECO 时间&安全管理平台\""
	},
	{
		"type": "body",
		"cms": "天融信安全管理系统-",
		"keyword": "title=\"天融信安全管理\""
	},
	{
		"type": "body",
		"cms": "锐捷 RG-DBS-",
		"keyword": "body=\"/css/impl-security.css\" || body=\"/dbaudit/authenticate\""
	},
	{
		"type": "body",
		"cms": "深信服防火墙类产品-",
		"keyword": "body=\"SANGFOR FW\""
	},
	{
		"type": "body",
		"cms": "天融信网络卫士过滤网关-",
		"keyword": "title=\"天融信网络卫士过滤网关\""
	},
	{
		"type": "body",
		"cms": "天融信网站监测与自动修复系统-",
		"keyword": "title=\"天融信网站监测与自动修复系统\""
	},
	{
		"type": "body",
		"cms": "天融信 TopAD-",
		"keyword": "title=\"天融信 TopAD\""
	},
	{
		"type": "body",
		"cms": "Apache-Forrest-",
		"keyword": "body=\"content=\"Apache Forrest\" || body=\"name=\"Forrest\""
	},
	{
		"type": "body",
		"cms": "Advantech-WebAccess-",
		"keyword": "body=\"/bw_templete1.dwt\" || body=\"/broadweb/WebAccessClientSetup.exe\" || body=\"/broadWeb/bwuconfig.asp\""
	},
	{
		"type": "body",
		"cms": "URP教务系统-",
		"keyword": "title=\"URP 综合教务系统\" || body=\"北京清元优软科技有限公司\""
	},
	{
		"type": "body",
		"cms": "H3C公司产品-",
		"keyword": "body=\"service@h3c.com\" || (body=\"Copyright\" && body=\"H3C Corporation\") || body=\"icg_helpScript.js\""
	},
	{
		"type": "body",
		"cms": "Huawei HG520 ADSL2+ Router-",
		"keyword": "title=\"Huawei HG520\""
	},
	{
		"type": "body",
		"cms": "Huawei B683V-",
		"keyword": "title=\"Huawei B683V\""
	},
	{
		"type": "body",
		"cms": "HUAWEI ESPACE 7910-",
		"keyword": "title=\"HUAWEI ESPACE 7910\""
	},
	{
		"type": "body",
		"cms": "Huawei HG630-",
		"keyword": "title=\"Huawei HG630\""
	},
	{
		"type": "body",
		"cms": "Huawei B683-",
		"keyword": "title=\"Huawei B683\""
	},
	{
		"type": "body",
		"cms": "华为 MCU-",
		"keyword": "body=\"McuR5-min.js\" || body=\"MCUType.js\" || title=\"huawei MCU\""
	},
	{
		"type": "body",
		"cms": "HUAWEI Inner Web-",
		"keyword": "title=\"HUAWEI Inner Web\" || body=\"hidden_frame.html\""
	},
	{
		"type": "body",
		"cms": "HUAWEI CSP-",
		"keyword": "title=\"HUAWEI CSP\""
	},
	{
		"type": "body",
		"cms": "华为 NetOpen-",
		"keyword": "body=\"/netopen/theme/css/inFrame.css\" || title=\"Huawei NetOpen System\""
	},
	{
		"type": "body",
		"cms": "校园卡管理系统-",
		"keyword": "body=\"Harbin synjones electronic\" || body=\"document.FormPostds.action=\"xxsearch.action\" || body=\"/shouyeziti.css\""
	},
	{
		"type": "body",
		"cms": "OBSERVA telcom-",
		"keyword": "title=\"OBSERVA\""
	},
	{
		"type": "body",
		"cms": "汉柏安全网关-",
		"keyword": "title=\"OPZOON - \""
	},
	{
		"type": "body",
		"cms": "b2evolution-",
		"keyword": "body=\"/powered-by-b2evolution-150t.gif\" || body=\"Powered by b2evolution\" || body=\"content=\"b2evolution\""
	},
	{
		"type": "body",
		"cms": "AvantFAX-",
		"keyword": "body=\"content=\"Web 2.0 HylaFAX\" || body=\"images/avantfax-big.png\""
	},
	{
		"type": "body",
		"cms": "Aurion-",
		"keyword": "body=\"<!-- Aurion Teal will be used as the login-time default\" || body=\"/aurion.js\""
	},
	{
		"type": "body",
		"cms": "Cisco-IP-Phone-",
		"keyword": "body=\"Cisco Unified Wireless IP Phone\""
	},
	{
		"type": "body",
		"cms": "Cisco-VPN-3000-Concentrator-",
		"keyword": "title=\"Cisco Systems, Inc. VPN 3000 Concentrator\""
	},
	{
		"type": "body",
		"cms": "BugTracker.NET-",
		"keyword": "body=\"href=\"btnet.css\" || body=\"valign=middle><a href=http://ifdefined.com/bugtrackernet.html>\" || body=\"<div class=logo>BugTracker.NET\""
	},
	{
		"type": "body",
		"cms": "BugFree-",
		"keyword": "body=\"id=\"logo\" alt=BugFree\" || body=\"class=\"loginBgImage\" alt=\"BugFree\" ||  title=\"BugFree\" || body=\"name=\"BugUserPWD\""
	},
	{
		"type": "body",
		"cms": "cPassMan-",
		"keyword": "title=\"Collaborative Passwords Manager\""
	},
	{
		"type": "body",
		"cms": "splunk-",
		"keyword": "body=\"Splunk.util.normalizeBoolean\""
	},
	{
		"type": "body",
		"cms": "DrugPak-",
		"keyword": "body=\"Powered by DrugPak\" || body=\"/dplimg/DPSTYLE.CSS\""
	},
	{
		"type": "body",
		"cms": "DMXReady-Portfolio-Manager-",
		"keyword": "body=\"/css/PortfolioManager/styles_display_page.css\" || body=\"rememberme_portfoliomanager\""
	},
	{
		"type": "body",
		"cms": "eGroupWare-",
		"keyword": "body=\"content=\"eGroupWare\""
	},
	{
		"type": "body",
		"cms": "eSyndiCat-",
		"keyword": "body=\"content=\"eSyndiCat\""
	},
	{
		"type": "body",
		"cms": "Epiware-",
		"keyword": "body=\"Epiware - Project and Document Management\""
	},
	{
		"type": "body",
		"cms": "eMeeting-Online-Dating-Software-",
		"keyword": "body=\"eMeeting Dating Software\" || body=\"/_eMeetingGlobals.js\""
	},
	{
		"type": "body",
		"cms": "FreeNAS-",
		"keyword": "body=\"title=\"Welcome to FreeNAS\" || body=\"/images/ui/freenas-logo.png\""
	},
	{
		"type": "body",
		"cms": "FestOS-",
		"keyword": "body=\"title=\"FestOS\" || body=\"css/festos.css\""
	},
	{
		"type": "body",
		"cms": "eTicket-",
		"keyword": "body=\"Powered by eTicket\" || body=\"<a href=\"http://www.eticketsupport.com\" target=\"_blank\">\" || body=\"/eticket/eticket.css\""
	},
	{
		"type": "body",
		"cms": "FileVista-",
		"keyword": "body=\"Welcome to FileVista\" || body=\"<a href=\"http://www.gleamtech.com/products/filevista/web-file-manager\""
	},
	{
		"type": "body",
		"cms": "Google-Talk-Chatback-",
		"keyword": "body=\"www.google.com/talk/service/\""
	},
	{
		"type": "body",
		"cms": "Flyspray-",
		"keyword": "body=\"Powered by Flyspray\""
	},
	{
		"type": "body",
		"cms": "HP-StorageWorks-Library-",
		"keyword": "title=\"HP StorageWorks\""
	},
	{
		"type": "body",
		"cms": "HostBill-",
		"keyword": "body=\"Powered by <a href=\"http://hostbillapp.com\" || body=\"<strong>HostBill\""
	},
	{
		"type": "body",
		"cms": "IBM-Cognos-",
		"keyword": "body=\"/cgi-bin/cognos.cgi\" || body=\"Cognos &#26159; International Business Machines Corp\""
	},
	{
		"type": "body",
		"cms": "iTop-",
		"keyword": "title=\"iTop Login\" || body=\"href=\"http://www.combodo.com/itop\""
	},
	{
		"type": "body",
		"cms": "Kayako-SupportSuite-",
		"keyword": "body=\"Powered By Kayako eSupport\" || body=\"Help Desk Software By Kayako eSupport\""
	},
	{
		"type": "body",
		"cms": "JXT-Consulting-",
		"keyword": "body=\"id=\"jxt-popup-wrapper\" || body=\"Powered by JXT Consulting\""
	},
	{
		"type": "body",
		"cms": "Fastly cdn-",
		"keyword": "body=\"fastcdn.org\""
	},
	{
		"type": "body",
		"cms": "JBoss_AS-",
		"keyword": "body=\"Manage this JBoss AS Instance\""
	},
	{
		"type": "body",
		"cms": "oracle_applicaton_server-",
		"keyword": "body=\"OraLightHeaderSub\""
	},
	{
		"type": "body",
		"cms": "Avaya-Aura-Utility-Server-",
		"keyword": "body=\"vmsTitle\">Avaya Aura&#8482;&nbsp;Utility Server\" || body=\"/webhelp/Base/Utility_toc.htm\""
	},
	{
		"type": "body",
		"cms": "DnP Firewall-",
		"keyword": "body=\"Powered by DnP Firewall\" || body=\"dnp_firewall_redirect\""
	},
	{
		"type": "body",
		"cms": "PaloAlto_Firewall-",
		"keyword": "body=\"Access to the web page you were trying to visit has been blocked in accordance with company policy\""
	},
	{
		"type": "body",
		"cms": "梭子鱼防火墙-",
		"keyword": "body=\"http://www.barracudanetworks.com?a=bsf_product\" class=\"transbutton\" && body=\"/cgi-mod/header_logo.cgi\""
	},
	{
		"type": "body",
		"cms": "IndusGuard_WAF-",
		"keyword": "title=\"IndusGuard WAF\" && body = \"wafportal/wafportal.nocache.js\""
	},
	{
		"type": "body",
		"cms": "网御WAF-",
		"keyword": "body = \"<div id=\"divLogin\">\" && title=\"网御WAF\""
	},
	{
		"type": "body",
		"cms": "NSFOCUS_WAF-",
		"keyword": "title=\"WAF NSFOCUS\" && body = \"/images/logo/nsfocus.png\""
	},
	{
		"type": "body",
		"cms": "斐讯Fortress-",
		"keyword": "title=\"斐讯Fortress防火墙\" && body=\"<meta name=\"author\" content=\"上海斐讯数据通信技术有限公司\" />\""
	},
	{
		"type": "body",
		"cms": "Sophos Web Appliance-",
		"keyword": "title=\"Sophos Web Appliance\" || body=\"resources/images/sophos_web.ico\" || body=\"url(resources/images/en/login_swa.jpg)\""
	},
	{
		"type": "body",
		"cms": "Barracuda-Spam-Firewall-",
		"keyword": "title=\"Barracuda Spam & Virus Firewall: Welcome\" || body=\"/barracuda.css\" || body=\"http://www.barracudanetworks.com?a=bsf_product\""
	},
	{
		"type": "body",
		"cms": "DnP-Firewall-",
		"keyword": "title=\"Forum Gateway - Powered by DnP Firewall\" || body=\"name=\"dnp_firewall_redirect\" ||  body=\"<form name=dnp_firewall\""
	},
	{
		"type": "body",
		"cms": "H3C-SecBlade-FireWall-",
		"keyword": "body=\"js/MulPlatAPI.js\""
	},
	{
		"type": "body",
		"cms": "锐捷NBR路由器-",
		"keyword": "body=\"free_nbr_login_form.png\""
	},
	{
		"type": "body",
		"cms": "mikrotik-",
		"keyword": "title=\"RouterOS\" && body=\"mikrotik\""
	},
	{
		"type": "body",
		"cms": "h3c路由器-",
		"keyword": "title=\"Web user login\" && body=\"nLanguageSupported\""
	},
	{
		"type": "body",
		"cms": "jcg无线路由器-",
		"keyword": "title=\"Wireless Router\" && body=\"http://www.jcgcn.com\""
	},
	{
		"type": "body",
		"cms": "Comcast_Business_Gateway-",
		"keyword": "body=\"Comcast Business Gateway\""
	},
	{
		"type": "body",
		"cms": "AirTiesRouter-",
		"keyword": "title=\"Airties\""
	},
	{
		"type": "body",
		"cms": "3COM NBX-",
		"keyword": "title=\"NBX NetSet\" || body=\"splashTitleIPTelephony\""
	},
	{
		"type": "body",
		"cms": "H3C ER2100n-",
		"keyword": "title=\"ER2100n系统管理\""
	},
	{
		"type": "body",
		"cms": "H3C ICG 1000-",
		"keyword": "title=\"ICG 1000系统管理\""
	},
	{
		"type": "body",
		"cms": "H3C AM8000-",
		"keyword": "title=\"AM8000\""
	},
	{
		"type": "body",
		"cms": "H3C ER8300G2-",
		"keyword": "title=\"ER8300G2系统管理\""
	},
	{
		"type": "body",
		"cms": "H3C ER3108GW-",
		"keyword": "title=\"ER3108GW系统管理\""
	},
	{
		"type": "body",
		"cms": "H3C ER6300-",
		"keyword": "title=\"ER6300系统管理\""
	},
	{
		"type": "body",
		"cms": "H3C ICG1000-",
		"keyword": "title=\"ICG1000系统管理\""
	},
	{
		"type": "body",
		"cms": "H3C ER3260G2-",
		"keyword": "title=\"ER3260G2系统管理\""
	},
	{
		"type": "body",
		"cms": "H3C ER3108G-",
		"keyword": "title=\"ER3108G系统管理\""
	},
	{
		"type": "body",
		"cms": "H3C ER2100-",
		"keyword": "title=\"ER2100系统管理\""
	},
	{
		"type": "body",
		"cms": "H3C ER3200-",
		"keyword": "title=\"ER3200系统管理\""
	},
	{
		"type": "body",
		"cms": "H3C ER8300-",
		"keyword": "title=\"ER8300系统管理\""
	},
	{
		"type": "body",
		"cms": "H3C ER5200G2-",
		"keyword": "title=\"ER5200G2系统管理\""
	},
	{
		"type": "body",
		"cms": "H3C ER6300G2-",
		"keyword": "title=\"ER6300G2系统管理\""
	},
	{
		"type": "body",
		"cms": "H3C ER2100V2-",
		"keyword": "title=\"ER2100V2系统管理\""
	},
	{
		"type": "body",
		"cms": "H3C ER3260-",
		"keyword": "title=\"ER3260系统管理\""
	},
	{
		"type": "body",
		"cms": "H3C ER3100-",
		"keyword": "title=\"ER3100系统管理\""
	},
	{
		"type": "body",
		"cms": "H3C ER5100-",
		"keyword": "title=\"ER5100系统管理\""
	},
	{
		"type": "body",
		"cms": "H3C ER5200-",
		"keyword": "title=\"ER5200系统管理\""
	},
	{
		"type": "body",
		"cms": "UBNT_UniFi系列路由-",
		"keyword": "title=\"UniFi\" && body=\"<div class=\"appGlobalHeader\">\""
	},
	{
		"type": "body",
		"cms": "AnyGate-",
		"keyword": "title=\"AnyGate\" || body=\"/anygate.php\""
	},
	{
		"type": "body",
		"cms": "Astaro-Security-Gateway-",
		"keyword": "body=\"wfe/asg/js/app_selector.js?t=\" || body=\"/doc/astaro-license.txt\" || body=\"/js/_variables_from_backend.js?t=\""
	},
	{
		"type": "body",
		"cms": "Aruba-Device-",
		"keyword": "body=\"/images/arubalogo.gif\" || (body=\"Copyright\" && body=\"Aruba Networks\")"
	},
	{
		"type": "body",
		"cms": "ARRIS-Touchstone-Router-",
		"keyword": "(body=\"Copyright\" && body=\"ARRIS Group\") || body=\"/arris_style.css\""
	},
	{
		"type": "body",
		"cms": "AP-Router-",
		"keyword": "title=\"AP Router New Generation\""
	},
	{
		"type": "body",
		"cms": "Belkin-Modem-",
		"keyword": "body=\"content=\"Belkin\""
	},
	{
		"type": "body",
		"cms": "Dell OpenManage Switch Administrator-",
		"keyword": "title=\"Dell OpenManage Switch Administrator\""
	},
	{
		"type": "body",
		"cms": "EDIMAX-",
		"keyword": "title=\"EDIMAX Technology\" || body=\"content=\"Edimax\""
	},
	{
		"type": "body",
		"cms": "eBuilding-Network-Controller-",
		"keyword": "title=\"eBuilding Web\""
	},
	{
		"type": "body",
		"cms": "ipTIME-Router-",
		"keyword": "title=\"networks - ipTIME\" || body=\"href=iptime.css\""
	},
	{
		"type": "body",
		"cms": "I-O-DATA-Router-",
		"keyword": "title=\"I-O DATA Wireless Broadband Router\""
	},
	{
		"type": "body",
		"cms": "phpshe-",
		"keyword": "body=\"Powered by phpshe\" || body=\"content=\"phpshe\""
	},
	{
		"type": "body",
		"cms": "ThinkSAAS-",
		"keyword": "body=\"/app/home/skins/default/style.css\""
	},
	{
		"type": "body",
		"cms": "e-tiller-",
		"keyword": "body=\"reader/view_abstract.aspx\""
	},
	{
		"type": "body",
		"cms": "DouPHP-",
		"keyword": "body=\"Powered by DouPHP\" || (body=\"controlBase\" && body=\"indexLeft\" && body=\"recommendProduct\")"
	},
	{
		"type": "body",
		"cms": "twcms-",
		"keyword": "body=\"/twcms/theme/\" && body=\"/css/global.css\""
	},
	{
		"type": "body",
		"cms": "SiteServer-",
		"keyword": "(body=\"Powered by\" && body=\"http://www.siteserver.cn\" && body=\"SiteServer CMS\") || title=\"Powered by SiteServer CMS\" || body=\"T_系统首页模板\" || (body=\"siteserver\" && body=\"sitefiles\")"
	},
	{
		"type": "body",
		"cms": "Joomla-",
		"keyword": "body=\"content=\"Joomla\" || (body=\"/media/system/js/core.js\" && body=\"/media/system/js/mootools-core.js\")"
	},
	{
		"type": "body",
		"cms": "kesionCMS-",
		"keyword": "body=\"/ks_inc/common.js\" || body=\"publish by KesionCMS\""
	},
	{
		"type": "body",
		"cms": "CMSTop-",
		"keyword": "body=\"/css/cmstop-common.css\" || body=\"/js/cmstop-common.js\" || body=\"cmstop-list-text.css\" || body=\"<a class=\"poweredby\" href=\"http://www.cmstop.com\"\""
	},
	{
		"type": "body",
		"cms": "ESPCMS-",
		"keyword": "title=\"Powered by ESPCMS\" || body=\"Powered by ESPCMS\" || (body=\"infolist_fff\" && body=\"/templates/default/style/tempates_div.css\")"
	},
	{
		"type": "body",
		"cms": "74cms-",
		"keyword": "body=\"content=\"74cms.com\" || body=\"content=\"骑士CMS\" || body=\"Powered by <a href=\"http://www.74cms.com/\"\" || (body=\"/templates/default/css/common.css\" && body=\"selectjobscategory\")"
	},
	{
		"type": "body",
		"cms": "Foosun-",
		"keyword": "body=\"Created by DotNetCMS\" || body=\"For Foosun\" || body=\"Powered by www.Foosun.net,Products:Foosun Content Manage system\""
	},
	{
		"type": "body",
		"cms": "PhpCMS-",
		"keyword": "(body=\"Powered by\" && body=\"http://www.phpcms.cn\") || body=\"content=\"Phpcms\" || body=\"Powered by Phpcms\" || body=\"data/config.js\" || body=\"/index.php?m=content&c=index&a=lists\" || body=\"/index.php?m=content&amp;c=index&amp;a=lists\""
	},
	{
		"type": "body",
		"cms": "DedeCMS-",
		"keyword": "body=\"Power by DedeCms\" || (body=\"Powered by\" && body=\"http://www.dedecms.com/\" && body=\"DedeCMS\") || body=\"/templets/default/style/dedecms.css\" || title=\"Powered by DedeCms\" "
	},
	{
		"type": "body",
		"cms": "ASPCMS-",
		"keyword": "title=\"Powered by ASPCMS\" || body=\"content=\"ASPCMS\" || body=\"/inc/AspCms_AdvJs.asp\""
	},
	{
		"type": "body",
		"cms": "MetInfo-",
		"keyword": "title=\"Powered by MetInfo\" || body=\"content=\"MetInfo\" || body=\"powered_by_metinfo\" || body=\"/images/css/metinfo.css\""
	},
	{
		"type": "body",
		"cms": "Npoint-",
		"keyword": "title=\"Powered by Npoint\""
	},
	{
		"type": "body",
		"cms": "捷点JCMS-",
		"keyword": "body=\"Publish By JCms2010\""
	},
	{
		"type": "body",
		"cms": "帝国EmpireCMS-",
		"keyword": "title=\"Powered by EmpireCMS\""
	},
	{
		"type": "body",
		"cms": "JEECMS-",
		"keyword": "title=\"Powered by JEECMS\" || (body=\"Powered by\" && body=\"http://www.jeecms.com\" && body=\"JEECMS\")"
	},
	{
		"type": "body",
		"cms": "IdeaCMS-",
		"keyword": "body=\"Powered By IdeaCMS\" || body=\"m_ctr32\""
	},
	{
		"type": "body",
		"cms": "TCCMS-",
		"keyword": "title=\"Power By TCCMS\" || (body=\"index.php?ac=link_more\" && body=\"index.php?ac=news_list\")"
	},
	{
		"type": "body",
		"cms": "webplus-",
		"keyword": "body=\"webplus\" && body=\"高校网站群管理平台\""
	},
	{
		"type": "body",
		"cms": "Dolibarr-",
		"keyword": "body=\"Dolibarr Development Team\""
	},
	{
		"type": "body",
		"cms": "Telerik Sitefinity-",
		"keyword": "body=\"Telerik.Web.UI.WebResource.axd\" || body=\"content=\"Sitefinity\""
	},
	{
		"type": "body",
		"cms": "PageAdmin-",
		"keyword": "body=\"content=\"PageAdmin CMS\"\" || body=\"/e/images/favicon.ico\""
	},
	{
		"type": "body",
		"cms": "sdcms-",
		"keyword": "title=\"powered by sdcms\" || (body=\"var webroot=\" && body=\"/js/sdcms.js\")"
	},
	{
		"type": "body",
		"cms": "EnterCRM-",
		"keyword": "body=\"EnterCRM\""
	},
	{
		"type": "body",
		"cms": "易普拉格科研管理系统-",
		"keyword": "body=\"lan12-jingbian-hong\" || body=\"科研管理系统，北京易普拉格科技\""
	},
	{
		"type": "body",
		"cms": "苏亚星校园管理系统-",
		"keyword": "body=\"/ws2004/Public/\""
	},
	{
		"type": "body",
		"cms": "trs_wcm-",
		"keyword": "body=\"/wcm/app/js\" || body=\"0;URL=/wcm\" || body=\"window.location.href = \"/wcm\";\" || (body=\"forum.trs.com.cn\" && body=\"wcm\") || body=\"/wcm\" target=\"_blank\">网站管理\" || body=\"/wcm\" target=\"_blank\">管理\""
	},
	{
		"type": "body",
		"cms": "we7-",
		"keyword": "body=\"/Widgets/WidgetCollection/\""
	},
	{
		"type": "body",
		"cms": "1024cms-",
		"keyword": "body=\"Powered by 1024 CMS\" || body=\"generator\" content=\"1024 CMS (c)\""
	},
	{
		"type": "body",
		"cms": "360webfacil_360WebManager-",
		"keyword": "(body=\"publico/template/\" && body=\"zonapie\") || body=\"360WebManager Software\""
	},
	{
		"type": "body",
		"cms": "6kbbs-",
		"keyword": "body=\"Powered by 6kbbs\" || body=\"generator\" content=\"6KBBS\""
	},
	{
		"type": "body",
		"cms": "Acidcat_CMS-",
		"keyword": "body=\"Start Acidcat CMS footer information\" || body=\"Powered by Acidcat CMS\""
	},
	{
		"type": "body",
		"cms": "bit-service-",
		"keyword": "body=\"bit-xxzs\" || body=\"xmlpzs/webissue.asp\""
	},
	{
		"type": "body",
		"cms": "云因网上书店-",
		"keyword": "body=\"main/building.cfm\" || body=\"href=\"../css/newscomm.css\""
	},
	{
		"type": "body",
		"cms": "MediaWiki-",
		"keyword": "body=\"generator\" content=\"MediaWiki\" || body=\"/wiki/images/6/64/Favicon.ico\" || body=\"Powered by MediaWiki\""
	},
	{
		"type": "body",
		"cms": "Typecho-",
		"keyword": "body=\"generator\" content=\"Typecho\" || (body=\"强力驱动\" && body=\"Typecho\")"
	},
	{
		"type": "body",
		"cms": "2z project-",
		"keyword": "body=\"Generator\" content=\"2z project\""
	},
	{
		"type": "body",
		"cms": "phpDocumentor-",
		"keyword": "body=\"Generated by phpDocumentor\""
	},
	{
		"type": "body",
		"cms": "微门户-",
		"keyword": "body=\"/tpl/Home/weimeng/common/css/\""
	},
	{
		"type": "body",
		"cms": "webEdition-",
		"keyword": "body=\"generator\" content=\"webEdition\""
	},
	{
		"type": "body",
		"cms": "orocrm-",
		"keyword": "body=\"/bundles/oroui/\""
	},
	{
		"type": "body",
		"cms": "创星伟业校园网群-",
		"keyword": "body=\"javascripts/float.js\" && body=\"vcxvcxv\""
	},
	{
		"type": "body",
		"cms": "BoyowCMS-",
		"keyword": "body=\"publish by BoyowCMS\""
	},
	{
		"type": "body",
		"cms": "正方教务管理系统-",
		"keyword": "body=\"style/base/jw.css\""
	},
	{
		"type": "body",
		"cms": "UFIDA_NC-",
		"keyword": "(body=\"UFIDA\" && body=\"logo/images/\") || body=\"logo/images/ufida_nc.png\""
	},
	{
		"type": "body",
		"cms": "phpweb-",
		"keyword": "body=\"PDV_PAGENAME\""
	},
	{
		"type": "body",
		"cms": "地平线CMS-",
		"keyword": "body=\"labelOppInforStyle\" || title=\"Powered by deep soon\" || (body=\"search_result.aspx\" && body=\"frmsearch\")"
	},
	{
		"type": "body",
		"cms": "HIMS酒店云计算服务-",
		"keyword": "(body=\"GB_ROOT_DIR\" && body=\"maincontent.css\") || body=\"HIMS酒店云计算服务\""
	},
	{
		"type": "body",
		"cms": "Tipask-",
		"keyword": "body=\"content=\"tipask\""
	},
	{
		"type": "body",
		"cms": "北创图书检索系统-",
		"keyword": "body=\"opac_two\""
	},
	{
		"type": "body",
		"cms": "微普外卖点餐系统-",
		"keyword": "body=\"Author\" content=\"微普外卖点餐系统\" || body=\"Powered By 点餐系统\" || body=\"userfiles/shoppics/\""
	},
	{
		"type": "body",
		"cms": "逐浪zoomla-",
		"keyword": "body=\"script src=\"http://code.zoomla.cn/\" || (body=\"NodePage.aspx\" && body=\"Item\") || body=\"/style/images/win8_symbol_140x140.png\""
	},
	{
		"type": "body",
		"cms": "北京清科锐华CEMIS-",
		"keyword": "body=\"/theme/2009/image\" && body=\"login.asp\""
	},
	{
		"type": "body",
		"cms": "asp168欧虎-",
		"keyword": "body=\"upload/moban/images/style.css\" || body=\"default.php?mod=article&do=detail&tid\""
	},
	{
		"type": "body",
		"cms": "擎天电子政务-",
		"keyword": "body=\"App_Themes/1/Style.css\" || body=\"window.location = \"homepages/index.aspx\" || body=\"homepages/content_page.aspx\""
	},
	{
		"type": "body",
		"cms": "北京阳光环球建站系统-",
		"keyword": "body=\"bigSortProduct.asp?bigid\""
	},
	{
		"type": "body",
		"cms": "MaticsoftSNS_动软分享社区-",
		"keyword": "body=\"MaticsoftSNS\" || (body=\"maticsoft\" && body=\"/Areas/SNS/\")"
	},
	{
		"type": "body",
		"cms": "FineCMS-",
		"keyword": "body=\"Powered by FineCMS\" || body=\"dayrui@gmail.com\" || body=\"Copyright\" content=\"FineCMS\""
	},
	{
		"type": "body",
		"cms": "Diferior-",
		"keyword": "body=\"Powered by Diferior\""
	},
	{
		"type": "body",
		"cms": "国家数字化学习资源中心系统-",
		"keyword": "title=\"页面加载中,请稍候\" && body=\"FrontEnd\""
	},
	{
		"type": "body",
		"cms": "某通用型政府cms-",
		"keyword": "body=\"/deptWebsiteAction.do\""
	},
	{
		"type": "body",
		"cms": "万户网络-",
		"keyword": "body=\"css/css_whir.css\""
	},
	{
		"type": "body",
		"cms": "rcms-",
		"keyword": "body=\"/r/cms/www/\" && body=\"jhtml\""
	},
	{
		"type": "body",
		"cms": "全国烟草系统-",
		"keyword": "body=\"ycportal/webpublish\""
	},
	{
		"type": "body",
		"cms": "O2OCMS-",
		"keyword": "body=\"/index.php/clasify/showone/gtitle/\""
	},
	{
		"type": "body",
		"cms": "一采通-",
		"keyword": "body=\"/custom/GroupNewsList.aspx?GroupId=\""
	},
	{
		"type": "body",
		"cms": "Dolphin-",
		"keyword": "body=\"bx_css_async\""
	},
	{
		"type": "body",
		"cms": "wecenter-",
		"keyword": "body=\"aw_template.js\" || body=\"WeCenter\""
	},
	{
		"type": "body",
		"cms": "phpvod-",
		"keyword": "body=\"Powered by PHPVOD\" || body=\"content=\"phpvod\""
	},
	{
		"type": "body",
		"cms": "08cms-",
		"keyword": "body=\"content=\"08CMS\" || body=\"typeof(_08cms)\""
	},
	{
		"type": "body",
		"cms": "tutucms-",
		"keyword": "body=\"content=\"TUTUCMS\" || body=\"Powered by TUTUCMS\" || body=\"TUTUCMS\"\""
	},
	{
		"type": "body",
		"cms": "八哥CMS-",
		"keyword": "body=\"content=\"BageCMS\""
	},
	{
		"type": "body",
		"cms": "mymps-",
		"keyword": "body=\"/css/mymps.css\" || title=\"mymps\" || body=\"content=\"mymps\""
	},
	{
		"type": "body",
		"cms": "IMGCms-",
		"keyword": "body=\"content=\"IMGCMS\" || body=\"Powered by IMGCMS\""
	},
	{
		"type": "body",
		"cms": "jieqi cms-",
		"keyword": "body=\"content=\"jieqi cms\" || title=\"jieqi cms\""
	},
	{
		"type": "body",
		"cms": "eadmin-",
		"keyword": "body=\"content=\"eAdmin\" || title=\"eadmin\""
	},
	{
		"type": "body",
		"cms": "opencms-",
		"keyword": "body=\"content=\"OpenCms\" || body=\"Powered by OpenCms\""
	},
	{
		"type": "body",
		"cms": "infoglue-",
		"keyword": "title=\"infoglue\" || body=\"infoglueBox.png\""
	},
	{
		"type": "body",
		"cms": "171cms-",
		"keyword": "body=\"content=\"171cms\" || title=\"171cms\""
	},
	{
		"type": "body",
		"cms": "doccms-",
		"keyword": "body=\"Power by DocCms\""
	},
	{
		"type": "body",
		"cms": "appcms-",
		"keyword": "body=\"Powerd by AppCMS\""
	},
	{
		"type": "body",
		"cms": "niucms-",
		"keyword": "body=\"content=\"NIUCMS\""
	},
	{
		"type": "body",
		"cms": "baocms-",
		"keyword": "body=\"content=\"BAOCMS\" || title=\"baocms\""
	},
	{
		"type": "body",
		"cms": "PublicCMS-",
		"keyword": "title=\"publiccms\""
	},
	{
		"type": "body",
		"cms": "JTBC(CMS)-",
		"keyword": "body=\"/js/jtbc.js\" || body=\"content=\"JTBC\""
	},
	{
		"type": "body",
		"cms": "易企CMS-",
		"keyword": "body=\"content=\"YiqiCMS\""
	},
	{
		"type": "body",
		"cms": "ZCMS-",
		"keyword": "body=\"_ZCMS_ShowNewMessage\" || body=\"zcms_skin\" || title=\"ZCMS泽元内容管理\""
	},
	{
		"type": "body",
		"cms": "科蚁CMS-",
		"keyword": "body=\"keyicms：keyicms\" || body=\"Powered by <a href=\"http://www.keyicms.com\""
	},
	{
		"type": "body",
		"cms": "苹果CMS-",
		"keyword": "body=\"maccms:voddaycount\""
	},
	{
		"type": "body",
		"cms": "大米CMS-",
		"keyword": "title=\"大米CMS-\" || body=\"content=\"damicms\" || body=\"content=\"大米CMS\""
	},
	{
		"type": "body",
		"cms": "phpmps-",
		"keyword": "body=\"Powered by Phpmps\" || body=\"templates/phpmps/style/index.css\""
	},
	{
		"type": "body",
		"cms": "25yi-",
		"keyword": "body=\"Powered by 25yi\" || body=\"css/25yi.css\""
	},
	{
		"type": "body",
		"cms": "kingcms-",
		"keyword": "title=\"kingcms\" || body=\"content=\"KingCMS\" || body=\"Powered by KingCMS\""
	},
	{
		"type": "body",
		"cms": "易点CMS-",
		"keyword": "body=\"DianCMS_SiteName\" || body=\"DianCMS_用户登陆引用\""
	},
	{
		"type": "body",
		"cms": "fengcms-",
		"keyword": "body=\"Powered by FengCms\" || body=\"content=\"FengCms\""
	},
	{
		"type": "body",
		"cms": "phpb2b-",
		"keyword": "body=\"Powered By PHPB2B\""
	},
	{
		"type": "body",
		"cms": "phpdisk-",
		"keyword": "body=\"Powered by PHPDisk\" || body=\"content=\"PHPDisk\""
	},
	{
		"type": "body",
		"cms": "EduSoho开源网络课堂-",
		"keyword": "title=\"edusoho\" || body=\"Powered by <a href=\"http://www.edusoho.com\" || body=\"Powered By EduSoho\""
	},
	{
		"type": "body",
		"cms": "phpok-",
		"keyword": "title=\"phpok\" || body=\"Powered By phpok.com\" || body=\"content=\"phpok\""
	},
	{
		"type": "body",
		"cms": "dtcms-",
		"keyword": "title=\"dtcms\" || body=\"content=\"动力启航,DTCMS\""
	},
	{
		"type": "body",
		"cms": "beecms-",
		"keyword": "(body=\"powerd by\" && body=\"BEESCMS\") || body=\"template/default/images/slides.min.jquery.js\""
	},
	{
		"type": "body",
		"cms": "ourphp-",
		"keyword": "body=\"content=\"OURPHP\" || body=\"Powered by ourphp\""
	},
	{
		"type": "body",
		"cms": "php云-",
		"keyword": "body=\"<div class=\"index_link_list_name\">\""
	},
	{
		"type": "body",
		"cms": "贷齐乐p2p-",
		"keyword": "body=\"/js/jPackageCss/jPackage.css\" || body=\"src=\"/js/jPackage\""
	},
	{
		"type": "body",
		"cms": "中企动力门户CMS-",
		"keyword": "body=\"中企动力提供技术支持\""
	},
	{
		"type": "body",
		"cms": "destoon-",
		"keyword": "body=\"<meta name=\"generator\" content=\"Destoon\" || body=\"destoon_moduleid\""
	},
	{
		"type": "body",
		"cms": "帝友P2P-",
		"keyword": "body=\"/js/diyou.js\" || body=\"src=\"/dyweb/dythemes\""
	},
	{
		"type": "body",
		"cms": "海洋CMS-",
		"keyword": "title=\"seacms\" || body=\"Powered by SeaCms\" || body=\"content=\"seacms\""
	},
	{
		"type": "body",
		"cms": "合正网站群内容管理系统-",
		"keyword": "body=\"Produced By\" && body=\"网站群内容管理系统\""
	},
	{
		"type": "body",
		"cms": "OpenSNS-",
		"keyword": "(body=\"powered by\" && body=\"opensns\") || body=\"content=\"OpenSNS\""
	},
	{
		"type": "body",
		"cms": "SEMcms-",
		"keyword": "body=\"semcms PHP\" || body=\"sc_mid_c_left_c sc_mid_left_bt\""
	},
	{
		"type": "body",
		"cms": "Yxcms-",
		"keyword": "body=\"/css/yxcms.css\" || body=\"content=\"Yxcms\""
	},
	{
		"type": "body",
		"cms": "NITC-",
		"keyword": "body=\"NITC Web Marketing Service\" || body=\"/images/nitc1.png\""
	},
	{
		"type": "body",
		"cms": "wuzhicms-",
		"keyword": "body=\"Powered by wuzhicms\" || body=\"content=\"wuzhicms\""
	},
	{
		"type": "body",
		"cms": "PHPMyWind-",
		"keyword": "body=\"phpMyWind.com All Rights Reserved\" || body=\"content=\"PHPMyWind\""
	},
	{
		"type": "body",
		"cms": "SiteEngine-",
		"keyword": "body=\"content=\"Boka SiteEngine\""
	},
	{
		"type": "body",
		"cms": "b2bbuilder-",
		"keyword": "body=\"content=\"B2Bbuilder\" || body=\"translateButtonId = \"B2Bbuilder\""
	},
	{
		"type": "body",
		"cms": "农友政务系统-",
		"keyword": "body=\"1207044504\""
	},
	{
		"type": "body",
		"cms": "dswjcms-",
		"keyword": "body=\"content=\"Dswjcms\" || body=\"Powered by Dswjcms\""
	},
	{
		"type": "body",
		"cms": "FoxPHP-",
		"keyword": "body=\"FoxPHPScroll\" || body=\"FoxPHP_ImList\" || body=\"content=\"FoxPHP\""
	},
	{
		"type": "body",
		"cms": "weiphp-",
		"keyword": "body=\"content=\"WeiPHP\" || body=\"/css/weiphp.css\""
	},
	{
		"type": "body",
		"cms": "iWebSNS-",
		"keyword": "body=\"/jooyea/images/sns_idea1.jpg\" || body=\"/jooyea/images/snslogo.gif\""
	},
	{
		"type": "body",
		"cms": "TurboCMS-",
		"keyword": "body=\"Powered by TurboCMS\" || body= \"/cmsapp/zxdcADD.jsp\" || body=\"/cmsapp/count/newstop_index.jsp?siteid=\""
	},
	{
		"type": "body",
		"cms": "MoMoCMS-",
		"keyword": "body=\"content=\"MoMoCMS\" || body=\"Powered BY MoMoCMS\""
	},
	{
		"type": "body",
		"cms": "Acidcat CMS-",
		"keyword": "body=\"Powered by Acidcat CMS\" || body=\"Start Acidcat CMS footer information\" || body=\"/css/admin_import.css\""
	},
	{
		"type": "body",
		"cms": "WP Plugin All-in-one-SEO-Pack-",
		"keyword": "body=\"<!-- /all in one seo pack -->\""
	},
	{
		"type": "body",
		"cms": "Aardvark Topsites-",
		"keyword": "body=\"Powered by\" && body=\"Aardvark Topsites\""
	},
	{
		"type": "body",
		"cms": "1024 CMS-",
		"keyword": "body=\"Powered by 1024 CMS\" || body=\"content=\"1024 CMS\""
	},
	{
		"type": "body",
		"cms": "68 Classifieds-",
		"keyword": "body=\"powered by\" && body=\"68 Classifieds\""
	},
	{
		"type": "body",
		"cms": "武汉弘智科技-",
		"keyword": "body=\"研发与技术支持：武汉弘智科技有限公司\""
	},
	{
		"type": "body",
		"cms": "北京金盘鹏图软件-",
		"keyword": "body=\"SpeakIntertScarch.aspx\""
	},
	{
		"type": "body",
		"cms": "育友软件-",
		"keyword": "body=\"http://www.yuysoft.com/\" && body=\"技术支持\""
	},
	{
		"type": "body",
		"cms": "STcms-",
		"keyword": "body=\"content=\"STCMS\" || body=\"DahongY<dahongy@gmail.com>\""
	},
	{
		"type": "body",
		"cms": "青果软件-",
		"keyword": "title=\"KINGOSOFT\" || body=\"SetKingoEncypt.jsp\" || body=\"/jkingo.js\""
	},
	{
		"type": "body",
		"cms": "DirCMS-",
		"keyword": "body=\"content=\"DirCMS\""
	},
	{
		"type": "body",
		"cms": "牛逼cms-",
		"keyword": "body=\"content=\"niubicms\""
	},
	{
		"type": "body",
		"cms": "南方数据-",
		"keyword": "body=\"/SouthidcKeFu.js\" || body=\"CONTENT=\"Copyright 2003-2015 - Southidc.net\" || body=\"/Southidcj2f.Js\""
	},
	{
		"type": "body",
		"cms": "yidacms-",
		"keyword": "body=\"yidacms.css\""
	},
	{
		"type": "body",
		"cms": "bluecms-",
		"keyword": "body=\"power by bcms\" || body=\"bcms_plugin\""
	},
	{
		"type": "body",
		"cms": "taocms-",
		"keyword": "body=\">taoCMS<\""
	},
	{
		"type": "body",
		"cms": "Tiki-wiki CMS-",
		"keyword": "body=\"jqueryTiki = new Object\""
	},
	{
		"type": "body",
		"cms": "lepton-cms-",
		"keyword": "body=\"content=\"LEPTON-CMS\" || body=\"Powered by LEPTON CMS\""
	},
	{
		"type": "body",
		"cms": "euse_study-",
		"keyword": "body=\"UserInfo/UserFP.aspx\""
	},
	{
		"type": "body",
		"cms": "沃科网异网同显系统-",
		"keyword": "body=\"沃科网\" || title=\"异网同显系统\""
	},
	{
		"type": "body",
		"cms": "Mixcall座席管理中心-",
		"keyword": "title=\"Mixcall座席管理中心\""
	},
	{
		"type": "body",
		"cms": "DuomiCms-",
		"keyword": "body=\"DuomiCms\" || title=\"Power by DuomiCms\""
	},
	{
		"type": "body",
		"cms": "ANECMS-",
		"keyword": "body=\"content=\"Erwin Aligam - ealigam@gmail.com\""
	},
	{
		"type": "body",
		"cms": "Ananyoo-CMS-",
		"keyword": "body=\"content=\"http://www.ananyoo.com\""
	},
	{
		"type": "body",
		"cms": "Amiro-CMS-",
		"keyword": "body=\"Powered by: Amiro CMS\" || body=\"-= Amiro.CMS (c) =-\""
	},
	{
		"type": "body",
		"cms": "AlumniServer-",
		"keyword": "body=\"AlumniServerProject.php\" || body=\"content=\"Alumni\""
	},
	{
		"type": "body",
		"cms": "AlstraSoft-EPay-Enterprise-",
		"keyword": "body=\"Powered by EPay Enterprise\" || body=\"/shop.htm?action=view\""
	},
	{
		"type": "body",
		"cms": "AlstraSoft-AskMe-",
		"keyword": "body=\"<a href=\"pass_recover.php\">\" || (body=\"Powered by\" && body=\"http://www.alstrasoft.com\")"
	},
	{
		"type": "body",
		"cms": "Artiphp-CMS-",
		"keyword": "body=\"copyright Artiphp\""
	},
	{
		"type": "body",
		"cms": "BIGACE-",
		"keyword": "body=\"content=\"BIGACE\" || body=\"Site is running BIGACE\""
	},
	{
		"type": "body",
		"cms": "Biromsoft-WebCam-",
		"keyword": "title=\"Biromsoft WebCam\""
	},
	{
		"type": "body",
		"cms": "BackBee-",
		"keyword": "body=\"<div id=\"bb5-site-wrapper\">\""
	},
	{
		"type": "body",
		"cms": "Auto-CMS-",
		"keyword": "body=\"Powered by Auto CMS\" || body=\"content=\"AutoCMS\""
	},
	{
		"type": "body",
		"cms": "STAR CMS-",
		"keyword": "body=\"content=\"STARCMS\" || body=\"<img alt=\"STAR CMS\""
	},
	{
		"type": "body",
		"cms": "Zotonic-",
		"keyword": "body=\"powered by: Zotonic\" || body=\"/lib/js/apps/zotonic-1.0\""
	},
	{
		"type": "body",
		"cms": "BloofoxCMS-",
		"keyword": "body=\"content=\"bloofoxCMS\" || body=\"Powered by <a href=\"http://www.bloofox.com\""
	},
	{
		"type": "body",
		"cms": "BlognPlus-",
		"keyword": "body=\"Powered by\" && body=\"href=\"http://www.blogn.org\""
	},
	{
		"type": "body",
		"cms": "bitweaver-",
		"keyword": "body=\"content=\"bitweaver\" || body=\"href=\"http://www.bitweaver.org\">Powered by\""
	},
	{
		"type": "body",
		"cms": "ClanSphere-",
		"keyword": "body=\"content=\"ClanSphere\" || body=\"index.php?mod=clansphere&amp;action=about\""
	},
	{
		"type": "body",
		"cms": "CitusCMS-",
		"keyword": "body=\"Powered by CitusCMS\" || body=\"<strong>CitusCMS</strong>\" || body=\"content=\"CitusCMS\""
	},
	{
		"type": "body",
		"cms": "CMS-WebManager-Pro-",
		"keyword": "body=\"content=\"Webmanager-pro\" || body=\"href=\"http://webmanager-pro.com\">Web.Manager\""
	},
	{
		"type": "body",
		"cms": "CMSQLite-",
		"keyword": "body=\"powered by CMSQLite\" || body=\"content=\"www.CMSQLite.net\""
	},
	{
		"type": "body",
		"cms": "CMSimple-",
		"keyword": "body=\"Powered by CMSimple.dk\" || body=\"content=\"CMSimple\""
	},
	{
		"type": "body",
		"cms": "CMScontrol-",
		"keyword": "body=\"content=\"CMScontrol\""
	},
	{
		"type": "body",
		"cms": "Claroline-",
		"keyword": "body=\"target=\"_blank\">Claroline</a>\" || body=\"http://www.claroline.net\" rel=\"Copyright\""
	},
	{
		"type": "body",
		"cms": "Car-Portal-",
		"keyword": "body=\"Powered by <a href=\"http://www.netartmedia.net/carsportal\" || body=\"class=\"bodyfontwhite\"><strong>&nbsp;Car Script\""
	},
	{
		"type": "body",
		"cms": "chillyCMS-",
		"keyword": "body=\"powered by <a href=\"http://FrozenPepper.de\""
	},
	{
		"type": "body",
		"cms": "BoonEx-Dolphin-",
		"keyword": "body=\"Powered by                    Dolphin - <a href=\"http://www.boonex.com/products/dolphin\""
	},
	{
		"type": "body",
		"cms": "SilverStripe-",
		"keyword": "body=\"content=\"SilverStripe\""
	},
	{
		"type": "body",
		"cms": "Campsite-",
		"keyword": "body=\"content=\"Campsite\""
	},
	{
		"type": "body",
		"cms": "ischoolsite-",
		"keyword": "body=\"Powered by <a href=\"http://www.ischoolsite.com\""
	},
	{
		"type": "body",
		"cms": "CafeEngine-",
		"keyword": "body=\"/CafeEngine/style.css\" || body=\"<a href=http://cafeengine.com>CafeEngine.com\""
	},
	{
		"type": "body",
		"cms": "BrowserCMS-",
		"keyword": "body=\"Powered by BrowserCMS\" || body=\"content=\"BrowserCMS\""
	},
	{
		"type": "body",
		"cms": "Contrexx-CMS-",
		"keyword": "body=\"powered by Contrexx\" || body=\"content=\"Contrexx\""
	},
	{
		"type": "body",
		"cms": "ContentXXL-",
		"keyword": "body=\"content=\"contentXXL\""
	},
	{
		"type": "body",
		"cms": "Contentteller-CMS-",
		"keyword": "body=\"content=\"Esselbach Contentteller CMS\""
	},
	{
		"type": "body",
		"cms": "Contao-",
		"keyword": "body=\"system/contao.css\""
	},
	{
		"type": "body",
		"cms": "CommonSpot-",
		"keyword": "body=\"content=\"CommonSpot\""
	},
	{
		"type": "body",
		"cms": "CruxCMS-",
		"keyword": "body=\"Created by CruxCMS\" || body=\"title=\"CruxCMS\" class=\"blank\""
	},
	{
		"type": "body",
		"cms": "锐商企业CMS-",
		"keyword": "body=\"href=\"/Writable/ClientImages/mycss.css\""
	},
	{
		"type": "body",
		"cms": "coWiki-",
		"keyword": "body=\"content=\"coWiki\" || body=\"<!-- Generated by coWiki\""
	},
	{
		"type": "body",
		"cms": "Coppermine-",
		"keyword": "body=\"<!--Coppermine Photo Gallery\""
	},
	{
		"type": "body",
		"cms": "DaDaBIK-",
		"keyword": "body=\"content=\"DaDaBIK\" || body=\"class=\"powered_by_dadabik\""
	},
	{
		"type": "body",
		"cms": "Custom-CMS-",
		"keyword": "body=\"content=\"CustomCMS\" || body=\"title=\"Powered by CCMS\""
	},
	{
		"type": "body",
		"cms": "DT-Centrepiece-",
		"keyword": "body=\"content=\"DT Centrepiece\" || body=\"Powered By DT Centrepiece\""
	},
	{
		"type": "body",
		"cms": "Edito-CMS-",
		"keyword": "body=\"content=\"edito\" || body=\"title=\"CMS\" href=\"http://www.edito.pl/\""
	},
	{
		"type": "body",
		"cms": "Echo-",
		"keyword": "body=\"powered by echo\" || body=\"/Echo2/echoweb/login\""
	},
	{
		"type": "body",
		"cms": "Ecomat-CMS-",
		"keyword": "body=\"content=\"ECOMAT CMS\""
	},
	{
		"type": "body",
		"cms": "EazyCMS-",
		"keyword": "body=\"powered by eazyCMS\" || body=\"<a class=\"actionlink\" href=\"http://www.eazyCMS.com\""
	},
	{
		"type": "body",
		"cms": "easyLink-Web-Solutions-",
		"keyword": "body=\"content=\"easyLink\""
	},
	{
		"type": "body",
		"cms": "EasyConsole-CMS-",
		"keyword": "body=\"Powered by EasyConsole CMS\" || body=\"Powered by <a href=\"http://www.easyconsole.com\""
	},
	{
		"type": "body",
		"cms": "DotCMS-",
		"keyword": "body=\"/dotAsset/\" || body=\"/index.dot\""
	},
	{
		"type": "body",
		"cms": "DBHcms-",
		"keyword": "body=\"powered by DBHcms\""
	},
	{
		"type": "body",
		"cms": "Donations-Cloud-",
		"keyword": "body=\"/donationscloud.css\""
	},
	{
		"type": "body",
		"cms": "Dokeos-",
		"keyword": "body=\"href=\"http://www.dokeos.com\" rel=\"Copyright\" || body=\"content=\"Dokeos\" || body=\"name=\"Generator\" content=\"Dokeos\""
	},
	{
		"type": "body",
		"cms": "Elxis-CMS-",
		"keyword": "body=\"content=\"Elxis\""
	},
	{
		"type": "body",
		"cms": "eFront-",
		"keyword": "body=\"<a href = \"http://www.efrontlearning.net\""
	},
	{
		"type": "body",
		"cms": "eSitesBuilder-",
		"keyword": "body=\"eSitesBuilder. All rights reserved\""
	},
	{
		"type": "body",
		"cms": "EPiServer-",
		"keyword": "body=\"content=\"EPiServer\" || body=\"/javascript/episerverscriptmanager.js\""
	},
	{
		"type": "body",
		"cms": "Energine-",
		"keyword": "body=\"scripts/Energine.js\" || body=\"Powered by <a href= \"http://energine.org/\" || body=\"stylesheets/energine.css\""
	},
	{
		"type": "body",
		"cms": "Gallery-",
		"keyword": "title=\"Gallery 3 Installer\" || body=\"/gallery/images/gallery.png\""
	},
	{
		"type": "body",
		"cms": "FrogCMS-",
		"keyword": "body=\"target=\"_blank\">Frog CMS\" || body=\"href=\"http://www.madebyfrog.com\">Frog CMS\""
	},
	{
		"type": "body",
		"cms": "Fossil-",
		"keyword": "body=\"<a href=\"http://fossil-scm.org\""
	},
	{
		"type": "body",
		"cms": "FCMS-",
		"keyword": "body=\"content=\"Ryan Haudenschilt\" || body=\"Powered by Family Connections\""
	},
	{
		"type": "body",
		"cms": "Fastpublish-CMS-",
		"keyword": "body=\"content=\"fastpublish\""
	},
	{
		"type": "body",
		"cms": "F3Site-",
		"keyword": "body=\"Powered by <a href=\"http://compmaster.prv.pl\""
	},
	{
		"type": "body",
		"cms": "Exponent-CMS-",
		"keyword": "body=\"content=\"Exponent Content Management System\" || body=\"Powered by Exponent CMS\""
	},
	{
		"type": "body",
		"cms": "E-Xoopport-",
		"keyword": "body=\"Powered by E-Xoopport\" || body=\"content=\"E-Xoopport\""
	},
	{
		"type": "body",
		"cms": "E-Manage-MySchool-",
		"keyword": "body=\"E-Manage All Rights Reserved MySchool Version\""
	},
	{
		"type": "body",
		"cms": "glFusion-",
		"keyword": "body=\"by <a href=\"http://www.glfusion.org/\""
	},
	{
		"type": "body",
		"cms": "GetSimple-",
		"keyword": "body=\"content=\"GetSimple\" || body=\"Powered by GetSimple\""
	},
	{
		"type": "body",
		"cms": "HESK-",
		"keyword": "body=\"hesk_javascript.js\" || body=\"hesk_style.css\" || body=\"Powered by <a href=\"http://www.hesk.com\" || body=\"Powered by <a href=\"https://www.hesk.com\""
	},
	{
		"type": "body",
		"cms": "GuppY-",
		"keyword": "body=\"content=\"GuppY\" || body=\"class=\"copyright\" href=\"http://www.freeguppy.org/\""
	},
	{
		"type": "body",
		"cms": "FluentNET-",
		"keyword": "body=\"content=\"Fluent\""
	},
	{
		"type": "body",
		"cms": "GeekLog-",
		"keyword": "body=\"Powered By <a href=\"http://www.geeklog.net/\""
	},
	{
		"type": "body",
		"cms": "Hycus-CMS-",
		"keyword": "body=\"content=\"Hycus\" || body=\"Powered By <a href=\"http://www.hycus.com\""
	},
	{
		"type": "body",
		"cms": "Hotaru-CMS-",
		"keyword": "body=\"content=\"Hotaru\""
	},
	{
		"type": "body",
		"cms": "HoloCMS-",
		"keyword": "body=\"Powered by HoloCMS\""
	},
	{
		"type": "body",
		"cms": "ImpressPages-CMS-",
		"keyword": "body=\"content=\"ImpressPages CMS\""
	},
	{
		"type": "body",
		"cms": "iGaming-CMS-",
		"keyword": "body=\"Powered by\" && body=\"http://www.igamingcms.com/\""
	},
	{
		"type": "body",
		"cms": "xoops-",
		"keyword": "body=\"include/xoops.js\""
	},
	{
		"type": "body",
		"cms": "Intraxxion-CMS-",
		"keyword": "body=\"content=\"Intraxxion\" || body=\"<!-- site built by Intraxxion\""
	},
	{
		"type": "body",
		"cms": "InterRed-",
		"keyword": "body=\"content=\"InterRed\" || body=\"Created with InterRed\""
	},
	{
		"type": "body",
		"cms": "Informatics-CMS-",
		"keyword": "body=\"content=\"Informatics\""
	},
	{
		"type": "body",
		"cms": "JagoanStore-",
		"keyword": "body=\"href=\"http://www.jagoanstore.com/\" target=\"_blank\">Toko Online\""
	},
	{
		"type": "body",
		"cms": "Kandidat-CMS-",
		"keyword": "body=\"content=\"Kandidat-CMS\""
	},
	{
		"type": "body",
		"cms": "Kajona-",
		"keyword": "body=\"content=\"Kajona\" || body=\"powered by Kajona\""
	},
	{
		"type": "body",
		"cms": "JGS-Portal-",
		"keyword": "body=\"Powered by <b>JGS-Portal Version\" || body=\"href=\"jgs_portal_box.php?id=\""
	},
	{
		"type": "body",
		"cms": "jCore-",
		"keyword": "body=\"JCORE_VERSION = \""
	},
	{
		"type": "body",
		"cms": "EdmWebVideo-",
		"keyword": "title=\"EdmWebVideo\""
	},
	{
		"type": "body",
		"cms": "edvr-",
		"keyword": "title=\"edvs/edvr\""
	},
	{
		"type": "body",
		"cms": "Polycom-",
		"keyword": "title=\"Polycom\" && body=\"kAllowDirectHTMLFileAccess\""
	},
	{
		"type": "body",
		"cms": "techbridge-",
		"keyword": "body=\"Sorry,you need to use IE brower\""
	},
	{
		"type": "body",
		"cms": "NETSurveillance-",
		"keyword": "title=\"NETSurveillance\""
	},
	{
		"type": "body",
		"cms": "nvdvr-",
		"keyword": "title=\"XWebPlay\""
	},
	{
		"type": "body",
		"cms": "DVR camera-",
		"keyword": "title=\"DVR WebClient\""
	},
	{
		"type": "body",
		"cms": "Macrec_DVR-",
		"keyword": "title=\"Macrec DVR\""
	},
	{
		"type": "body",
		"cms": "OnSSI_Video_Clients-",
		"keyword": "title=\"OnSSI Video Clients\" || body=\"x-value=\"On-Net Surveillance Systems Inc.\"\""
	},
	{
		"type": "body",
		"cms": "Linksys_SPA_Configuration -",
		"keyword": "title=\"Linksys SPA Configuration\""
	},
	{
		"type": "body",
		"cms": "eagleeyescctv-",
		"keyword": "body=\"IP Surveillance for Your Life\" || body=\"/nobody/loginDevice.js\""
	},
	{
		"type": "body",
		"cms": "dasannetworks-",
		"keyword": "body=\"clear_cookie(\"login\");\""
	},
	{
		"type": "body",
		"cms": "海康威视iVMS-",
		"keyword": "body=\"g_szCacheTime\" && body=\"iVMS\""
	},
	{
		"type": "body",
		"cms": "佳能网络摄像头(Canon Network Cameras)-",
		"keyword": "body=\"/viewer/live/en/live.html\""
	},
	{
		"type": "body",
		"cms": "NetDvrV3-",
		"keyword": "body=\"objLvrForNoIE\""
	},
	{
		"type": "body",
		"cms": "SIEMENS IP Cameras-",
		"keyword": "title=\"SIEMENS IP Camera\""
	},
	{
		"type": "body",
		"cms": "VideoIQ Camera-",
		"keyword": "title=\"VideoIQ Camera Login\""
	},
	{
		"type": "body",
		"cms": "Honeywell IP-Camera-",
		"keyword": "title=\"Honeywell IP-Camera\""
	},
	{
		"type": "body",
		"cms": "sony摄像头-",
		"keyword": "title=\"Sony Network Camera\" || body=\"inquiry.cgi?inqjs=system&inqjs=camera\""
	},
	{
		"type": "body",
		"cms": "AJA-Video-Converter-",
		"keyword": "body=\"eParamID_SWVersion\""
	},
	{
		"type": "body",
		"cms": "ACTi-",
		"keyword": "title=\"Web Configurator\" || body=\"ACTi Corporation All Rights Reserved\""
	},
	{
		"type": "body",
		"cms": "Samsung DVR-",
		"keyword": "title=\"Samsung DVR\""
	},
	{
		"type": "body",
		"cms": "Vicworl-",
		"keyword": "body=\"Powered by Vicworl\" || body=\"content=\"Vicworl\" || body=\"vindex_right_d\""
	},
	{
		"type": "body",
		"cms": "AVCON6-",
		"keyword": "body=\"filename=AVCON6Setup.exe\" || title=\"AVCON6系统管理平台\"  || body=\"language_dispose.action\""
	},
	{
		"type": "body",
		"cms": "Axis-Network-Camera-",
		"keyword": "title=\"AXIS Video Server\" || body=\"/incl/trash.shtml\""
	},
	{
		"type": "body",
		"cms": "Panasonic Network Camera-",
		"keyword": "body=\"MultiCameraFrame?Mode=Motion&Language\""
	},
	{
		"type": "body",
		"cms": "BlueNet-Video-",
		"keyword": "body=\"/cgi-bin/client_execute.cgi?tUD=0\" || title=\"BlueNet Video Viewer Version\""
	},
	{
		"type": "body",
		"cms": "ClipBucket-",
		"keyword": "body=\"content=\"ClipBucket\" || body=\"<!-- ClipBucket\" || body=\"<!-- Forged by ClipBucket\" || body=\"href=\"http://clip-bucket.com/\">ClipBucket\""
	},
	{
		"type": "body",
		"cms": "ZoneMinder-",
		"keyword": "body=\"ZoneMinder Login\""
	},
	{
		"type": "body",
		"cms": "DVR-WebClient-",
		"keyword": "body=\"259F9FDF-97EA-4C59-B957-5160CAB6884E\" || title=\"DVR-WebClient\""
	},
	{
		"type": "body",
		"cms": "D-Link-Network-Camera-",
		"keyword": "body=\"DCS-950G\".toLowerCase()\" || title=\"DCS-5300\""
	},
	{
		"type": "body",
		"cms": "DiBos-",
		"keyword": "title=\"DiBos - Login\" || body=\"style/bovisnt.css\""
	},
	{
		"type": "body",
		"cms": "Evo-Cam-",
		"keyword": "body=\"value=\"evocam.jar\" || body=\"<applet archive=\"evocam.jar\""
	},
	{
		"type": "body",
		"cms": "Intellinet-IP-Camera-",
		"keyword": "body=\"Copyright &copy;  INTELLINET NETWORK SOLUTIONS\" || body=\"http://www.intellinet-network.com/driver/NetCam.exe\""
	},
	{
		"type": "body",
		"cms": "IQeye-Netcam-",
		"keyword": "title=\"IQEYE: Live Images\" || body=\"content=\"Brian Lau, IQinVision\" || body=\"loc = \"iqeyevid.html\""
	},
	{
		"type": "body",
		"cms": "phpwind-",
		"keyword": "title=\"Powered by phpwind\" || body=\"content=\"phpwind\""
	},
	{
		"type": "body",
		"cms": "discuz-",
		"keyword": "title=\"Powered by Discuz\" || body=\"content=\"Discuz\" || (body=\"discuz_uid\" && body=\"portal.php?mod=view\") || body=\"Powered by <strong><a href=\"http://www.discuz.net\""
	},
	{
		"type": "body",
		"cms": "6kbbs-",
		"keyword": "body=\"Powered by 6kbbs\" || body=\"generator\" content=\"6KBBS\""
	},
	{
		"type": "body",
		"cms": "IP.Board-",
		"keyword": "body=\"ipb.vars\""
	},
	{
		"type": "body",
		"cms": "ThinkOX-",
		"keyword": "body=\"Powered By ThinkOX\" || title=\"ThinkOX\""
	},
	{
		"type": "body",
		"cms": "bbPress-",
		"keyword": "body=\"<!-- If you like showing off the fact that your server rocks -->\" || body=\"is proudly powered by <a href=\"http://bbpress.org\""
	},
	{
		"type": "body",
		"cms": "BlogEngine_NET-",
		"keyword": "body=\"pics/blogengine.ico\" || (body=\"Powered by\" && body=\"http://www.dotnetblogengine.net\")"
	},
	{
		"type": "body",
		"cms": "boastMachine-",
		"keyword": "body=\"powered by boastMachine\" || body=\"Powered by <a href=\"http://boastology.com\""
	},
	{
		"type": "body",
		"cms": "BrewBlogger-",
		"keyword": "body=\"developed by <a href=\"http://www.zkdigital.com\""
	},
	{
		"type": "body",
		"cms": "Dotclear-",
		"keyword": "body=\"Powered by <a href=\"http://dotclear.org/\""
	},
	{
		"type": "body",
		"cms": "DokuWiki-",
		"keyword": "body=\"powered by DokuWiki\" || body=\"content=\"DokuWiki\" || body=\"<div id=\"dokuwiki\""
	},
	{
		"type": "body",
		"cms": "DeluxeBB-",
		"keyword": "body=\"content=\"powered by DeluxeBB\""
	},
	{
		"type": "body",
		"cms": "esoTalk-",
		"keyword": "body=\"generated by esoTalk\" || body=\"Powered by esoTalk\" || body=\"/js/esotalk.js\""
	},
	{
		"type": "body",
		"cms": "Hiki-",
		"keyword": "body=\"content=\"Hiki\" || body=\"/hiki_base.css\" || body=\"by <a href=\"http://hikiwiki.org/\""
	},
	{
		"type": "body",
		"cms": "Gossamer-Forum-",
		"keyword": "body=\"href=\"gforum.cgi?username=\" || title=\"Gossamer Forum\""
	},
	{
		"type": "body",
		"cms": "Forest-Blog-",
		"keyword": "title=\"Forest Blog\""
	},
	{
		"type": "body",
		"cms": "FluxBB-",
		"keyword": "body=\"Powered by <a href=\"http://fluxbb.org/\""
	},
	{
		"type": "body",
		"cms": "Kampyle-",
		"keyword": "body=\"http://cf.kampyle.com/k_button.js\" || body=\"Start Kampyle Feedback Form Button\""
	},
	{
		"type": "body",
		"cms": "KaiBB-",
		"keyword": "body=\"Powered by KaiBB\" || body=\"content=\"Forum powered by KaiBB\""
	},
	{
		"type": "body",
		"cms": "fangmail-",
		"keyword": "body=\"/fangmail/default/css/em_css.css\""
	},
	{
		"type": "body",
		"cms": "MDaemon-",
		"keyword": "body=\"/WorldClient.dll?View=Main\""
	},
	{
		"type": "body",
		"cms": "网易企业邮箱-",
		"keyword": "body=\"frmvalidator\" && title=\"邮箱用户登录\""
	},
	{
		"type": "body",
		"cms": "TurboMail-",
		"keyword": "body=\"Powered by TurboMail\" || body=\"wzcon1 clearfix\" || title=\"TurboMail邮件系统\""
	},
	{
		"type": "body",
		"cms": "万网企业云邮箱-",
		"keyword": "body=\"static.mxhichina.com/images/favicon.ico\""
	},
	{
		"type": "body",
		"cms": "bxemail-",
		"keyword": "title=\"百讯安全邮件系统\" || title=\"百姓邮局\" || body=\"请输入正确的电子邮件地址，如：abc@bxemail.com\""
	},
	{
		"type": "body",
		"cms": "Coremail-",
		"keyword": "title=\"/coremail/common/assets\" || title=\"Coremail邮件系统\""
	},
	{
		"type": "body",
		"cms": "Lotus-",
		"keyword": "title=\"IBM Lotus iNotes Login\" || body=\"iwaredir.nsf\""
	},
	{
		"type": "body",
		"cms": "mirapoint-",
		"keyword": "body=\"/wm/mail/login.html\""
	},
	{
		"type": "body",
		"cms": "U-Mail-",
		"keyword": "body=\"<BODY LINK=\"White\" VLINK=\"White\" ALINK=\"White\">\""
	},
	{
		"type": "body",
		"cms": "Spammark邮件信息安全网关-",
		"keyword": "title=\"Spammark邮件信息安全网关\" || body=\"/cgi-bin/spammark?empty=1\""
	},
	{
		"type": "body",
		"cms": "科信邮件系统-",
		"keyword": "body=\"/systemfunction.pack.js\" || body=\"lo_computername\""
	},
	{
		"type": "body",
		"cms": "winwebmail-",
		"keyword": "title=\"winwebmail\" || body=\"WinWebMail Server\"  || body=\"images/owin.css\""
	},
	{
		"type": "body",
		"cms": "泰信TMailer邮件系统-",
		"keyword": "title=\"Tmailer\" || body=\"content=\"Tmailer\" || body=\"href=\"/tmailer/img/logo/favicon.ico\""
	},
	{
		"type": "body",
		"cms": "richmail-",
		"keyword": "title=\"Richmail\" || body=\"/resource/se/lang/se/mail_zh_CN.js\" || body=\"content=\"Richmail\""
	},
	{
		"type": "body",
		"cms": "iGENUS邮件系统-",
		"keyword": "body=\"Copyright by<A HREF=\"http://www.igenus.org\" || title=\"iGENUS webmail\""
	},
	{
		"type": "body",
		"cms": "金笛邮件系统-",
		"keyword": "body=\"/jdwm/cgi/login.cgi?login\""
	},
	{
		"type": "body",
		"cms": "迈捷邮件系统(MagicMail)-",
		"keyword": "body=\"/aboutus/magicmail.gif\""
	},
	{
		"type": "body",
		"cms": "Atmail-WebMail-",
		"keyword": "body=\"Powered by Atmail\" || body=\"/index.php/mail/auth/processlogin\" || body=\"<input id=\"Mailserverinput\""
	},
	{
		"type": "body",
		"cms": "FormMail-",
		"keyword": "body=\"/FormMail.pl\" || body=\"href=\"http://www.worldwidemart.com/scripts/formmail.shtml\""
	},
	{
		"type": "body",
		"cms": "同城多用户商城-",
		"keyword": "body=\"style_chaoshi\""
	},
	{
		"type": "body",
		"cms": "iWebShop-",
		"keyword": "body=\"/runtime/default/systemjs\""
	},
	{
		"type": "body",
		"cms": "1und1-",
		"keyword": "body=\"/shop/catalog/browse?sessid=\""
	},
	{
		"type": "body",
		"cms": "cart_engine-",
		"keyword": "body=\"skins/_common/jscripts.css\""
	},
	{
		"type": "body",
		"cms": "Magento-",
		"keyword": "(body=\"/skin/frontend/\" && body=\"BLANK_IMG\") || body=\"Magento, Varien, E-commerce\""
	},
	{
		"type": "body",
		"cms": "OpenCart-",
		"keyword": "body=\"Powered By OpenCart\" || body=\"catalog/view/theme\""
	},
	{
		"type": "body",
		"cms": "hishop-",
		"keyword": "body=\"hishop.plugins.openid\" || body=\"Hishop development team\""
	},
	{
		"type": "body",
		"cms": "Maticsoft_Shop_动软商城-",
		"keyword": "body=\"Maticsoft Shop\" || (body=\"maticsoft\" && body=\"/Areas/Shop/\")"
	},
	{
		"type": "body",
		"cms": "hikashop-",
		"keyword": "body=\"/media/com_hikashop/css/\""
	},
	{
		"type": "body",
		"cms": "tp-shop-",
		"keyword": "body=\"mn-c-top\""
	},
	{
		"type": "body",
		"cms": " 海盗云商(Haidao)-",
		"keyword": "body=\"haidao.web.general.js\""
	},
	{
		"type": "body",
		"cms": "shopbuilder-",
		"keyword": "body=\"content=\"ShopBuilder\" || body=\"Powered by ShopBuilder\" || body=\"ShopBuilder版权所有\""
	},
	{
		"type": "body",
		"cms": "v5shop-",
		"keyword": "title=\"v5shop\" || body=\"content=\"V5shop\" || body=\"Powered by V5Shop\""
	},
	{
		"type": "body",
		"cms": "shopnc-",
		"keyword": "body=\"Powered by ShopNC\" || body=\"Copyright 2007-2014 ShopNC Inc\" || body=\"content=\"ShopNC\""
	},
	{
		"type": "body",
		"cms": "shopex-",
		"keyword": "body=\"content=\"ShopEx\" || body=\"@author litie[aita]shopex.cn\""
	},
	{
		"type": "body",
		"cms": "dbshop-",
		"keyword": "body=\"content=\"dbshop\""
	},
	{
		"type": "body",
		"cms": "任我行电商-",
		"keyword": "body=\"content=\"366EC\""
	},
	{
		"type": "body",
		"cms": "CuuMall-",
		"keyword": "body=\"Power by CuuMall\""
	},
	{
		"type": "body",
		"cms": "javashop-",
		"keyword": "body=\"易族智汇javashop\" || body=\"javashop微信公众号\" || body=\"content=\"JavaShop\""
	},
	{
		"type": "body",
		"cms": "TPshop-",
		"keyword": "body=\"/index.php/Mobile/Index/index.html\" || body=\">TPshop开源商城<\""
	},
	{
		"type": "body",
		"cms": "MvMmall-",
		"keyword": "body=\"content=\"MvMmall\""
	},
	{
		"type": "body",
		"cms": "AirvaeCommerce-",
		"keyword": "body=\"E-Commerce Shopping Cart Software\""
	},
	{
		"type": "body",
		"cms": "AiCart-",
		"keyword": "body=\"APP_authenticate\""
	},
	{
		"type": "body",
		"cms": "MallBuilder-",
		"keyword": "body=\"content=\"MallBuilder\" || body=\"Powered by MallBuilder\""
	},
	{
		"type": "body",
		"cms": "e-junkie-",
		"keyword": "body=\"function EJEJC_lc\""
	},
	{
		"type": "body",
		"cms": "Allomani-",
		"keyword": "body=\"content=\"Allomani\" || body=\"Programmed By Allomani\""
	},
	{
		"type": "body",
		"cms": "ASPilot-Cart-",
		"keyword": "body=\"content=\"Pilot Cart\" || body=\"/pilot_css_default.css\""
	},
	{
		"type": "body",
		"cms": "Axous-",
		"keyword": "body=\"content=\"Axous\" || body=\"title=\"Axous Shareware Shop\""
	},
	{
		"type": "body",
		"cms": "CaupoShop-Classic-",
		"keyword": "body=\"Powered by CaupoShop\" || body=\"<!-- CaupoShop Classic\" || body=\"<a href=\"http://www.caupo.net\" target=\"_blank\">CaupoNet\""
	},
	{
		"type": "body",
		"cms": "PretsaShop-",
		"keyword": "body=\"content=\"PrestaShop\"\""
	},
	{
		"type": "body",
		"cms": "ComersusCart-",
		"keyword": "body=\"CONTENT=\"Powered by Comersus\" || body=\"href=\"comersus_showCart.asp\""
	},
	{
		"type": "body",
		"cms": "Foxycart-",
		"keyword": "body=\"<script src=\"//cdn.foxycart.com\""
	},
	{
		"type": "body",
		"cms": "DV-Cart-",
		"keyword": "body=\"class=\"KT_tngtable\""
	},
	{
		"type": "body",
		"cms": "EarlyImpact-ProductCart-",
		"keyword": "body=\"fpassword.asp?redirectUrl=&frURL=Custva.asp\""
	},
	{
		"type": "body",
		"cms": "Escenic-",
		"keyword": "body=\"content=\"Escenic\" || body=\"<!-- Start Escenic Analysis Engine client script -->\""
	},
	{
		"type": "body",
		"cms": "ICEshop-",
		"keyword": "body=\"Powered by ICEshop\" || body=\"<div id=\"iceshop\">\""
	},
	{
		"type": "body",
		"cms": "Interspire-Shopping-Cart-",
		"keyword": "body=\"content=\"Interspire Shopping Cart\" || body=\"class=\"PoweredBy\">Interspire Shopping Cart\""
	},
	{
		"type": "body",
		"cms": "iScripts-MultiCart-",
		"keyword": "body=\"Powered by <a href=\"http://iscripts.com/multicart\""
	},
	{
		"type": "body",
		"cms": "华天动力OA(OA8000)-",
		"keyword": "body=\"/OAapp/WebObjects/OAapp.woa\""
	},
	{
		"type": "body",
		"cms": "通达OA-",
		"keyword": "body=\"<link rel=\"shortcut icon\" href=\"/images/tongda.ico\" />\" || (body=\"OA提示：不能登录OA\" && body=\"紧急通知：今日10点停电\") || body=\"Office Anywhere 2013\"|| body = \"<a href='http://www.tongda2000.com/' target='_black'>通达官网</a></div>\""
	},
	{
		"type": "body",
		"cms": "OA(a8/seeyon/ufida)-",
		"keyword": "body=\"/seeyon/USER-DATA/IMAGES/LOGIN/login.gif\""
	},
	{
		"type": "body",
		"cms": "yongyoufe-",
		"keyword": "title=\"FE协作\" || (body=\"V_show\" && body=\"V_hedden\")"
	},
	{
		"type": "body",
		"cms": "pmway_E4_crm-",
		"keyword": "title=\"E4\" && title=\"CRM\""
	},
	{
		"type": "body",
		"cms": "Dolibarr-",
		"keyword": "body=\"Dolibarr Development Team\""
	},
	{
		"type": "body",
		"cms": "PHPOA-",
		"keyword": "body=\"admin_img/msg_bg.png\""
	},
	{
		"type": "body",
		"cms": "78oa-",
		"keyword": "body=\"/resource/javascript/system/runtime.min.js\" || body=\"license.78oa.com\" || title=\"78oa\"||body=\"src=\"/module/index.php\""
	},
	{
		"type": "body",
		"cms": "WishOA-",
		"keyword": "body=\"WishOA_WebPlugin.js\""
	},
	{
		"type": "body",
		"cms": "金和协同管理平台-",
		"keyword": "title=\"金和协同管理平台\""
	},
	{
		"type": "body",
		"cms": "Lotus-",
		"keyword": "title=\"IBM Lotus iNotes Login\" || body=\"iwaredir.nsf\""
	},
	{
		"type": "body",
		"cms": "OA企业智能办公自动化系统-",
		"keyword": "body=\"input name=\"S1\" type=\"image\"\" && body=\"count/mystat.asp\""
	},
	{
		"type": "body",
		"cms": "ecwapoa-",
		"keyword": "body=\"ecwapoa\""
	},
	{
		"type": "body",
		"cms": "ezOFFICE-",
		"keyword": "title=\"Wanhu ezOFFICE\" || body=\"EZOFFICEUSERNAME\" ||title=\"万户OA\" || body=\"whirRootPath\" || body=\"/defaultroot/js/cookie.js\""
	},
	{
		"type": "body",
		"cms": "任我行CRM-",
		"keyword": "title=\"任我行CRM\" || body=\"CRM_LASTLOGINUSERKEY\""
	},
	{
		"type": "body",
		"cms": "信达OA-",
		"keyword": "body=\"http://www.xdoa.cn</a>\" || body=\"北京创信达科技有限公司\""
	},
	{
		"type": "body",
		"cms": "协众OA-",
		"keyword": "body= \"Powered by 协众OA\" || body=\"admin@cnoa.cn\" || body=\"Powered by CNOA.CN\""
	},
	{
		"type": "body",
		"cms": "soffice-",
		"keyword": "title=\"OA办公管理平台\""
	},
	{
		"type": "body",
		"cms": "海天OA-",
		"keyword": "body=\"HTVOS.js\""
	},
	{
		"type": "body",
		"cms": "泛微OA-java",
		"keyword": "body=\"/js/jquery/jquery_wev8.js\"||body=\"/login/Login.jsp?logintype=1\""
	},
	{
		"type": "body",
		"cms": "中望OA-",
		"keyword": "body=\"/app_qjuserinfo/qjuserinfoadd.jsp\" || body=\"/IMAGES/default/first/xtoa_logo.png\""
	},
	{
		"type": "body",
		"cms": "睿博士云办公系统-",
		"keyword": "body=\"/studentSign/toLogin.di\" || body=\"/user/toUpdatePasswordPage.di\""
	},
	{
		"type": "body",
		"cms": "一米OA-",
		"keyword": "body=\"/yimioa.apk\""
	},
	{
		"type": "body",
		"cms": "泛普建筑工程施工OA-",
		"keyword": "body=\"/dwr/interface/LoginService.js\""
	},
	{
		"type": "body",
		"cms": "正方OA-",
		"keyword": "body=\"zfoausername\""
	},
	{
		"type": "body",
		"cms": "希尔OA-",
		"keyword": "body=\"/heeroa/login.do\""
	},
	{
		"type": "body",
		"cms": "用友致远oa-",
		"keyword": "body=\"/seeyon/USER-DATA/IMAGES/LOGIN/login.gif\" || title=\"用友致远A\" || body=\"/yyoa/\" || body=\"/seeyon/common/all-min.js\""
	},
	{
		"type": "body",
		"cms": "WordPress-php",
		"keyword": "body=\"/wp-login.php?\"||body=\"wp-user\""
	},
	{
		"type": "body",
		"cms": "宝塔面板-python",
		"keyword": "body=\"<title>安全入口校验失败</title>\" || body=\"https://www.bt.cn/bbs/thread-18367-1-1.html\""
	},
	{
		"type": "body",
		"cms": "Emlog-PHP",
		"keyword": "body=\"/include/lib/js/common_tpl.js\"&&body=\"content/templates\""
	}
]
`)
