package tls

import (
	"LiteKube/pkg/common"
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"
)

type TLSInfo struct {
	CACert     string `json:"CACert"`
	ClientCert string `json:"ClientCert"`
	ClientKey  string `json:"ClientKey"`
	Port       int    `json:"Port"`
}

func (tlsInfo TLSInfo) HTMLBR() TLSInfo {
	return TLSInfo{
		CACert:     strings.Replace(tlsInfo.CACert, "\n", "<br>", -1),
		ClientCert: strings.Replace(tlsInfo.ClientCert, "\n", "<br>", -1),
		ClientKey:  strings.Replace(tlsInfo.ClientKey, "\n", "<br>", -1),
		Port:       tlsInfo.Port,
	}
}

func TLSReturnHTML(info TLSInfo) string {
	buffer := new(bytes.Buffer)
	temp, err := template.New("").Parse(tlsHtmlTemplate)
	if err != nil {
		return fmt.Sprintf("We meet some errors: %s", err)
	}

	err = temp.Execute(buffer, info.HTMLBR())
	if err != nil {
		return fmt.Sprintf("We meet some errors: %s", err)
	}

	return buffer.String()
}

func TLSReturnJson(info TLSInfo) string {
	jsonBytes, err := json.Marshal(info)
	if err != nil {
		return common.ErrorJson("fail to Marshal json")
	}
	return string(jsonBytes)
}

func TLSReturn(info TLSInfo, isRaw bool) string {
	if isRaw {
		return TLSReturnHTML(info)
	} else {
		return TLSReturnJson(info)
	}
}

var tlsHtmlTemplate string = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8"> 
		<title>LiteKube Certificate Page</title> 
		<style>
			p.double {border-style:double;}
			p.ridge {border-style:ridge;}
		</style>
	</head>

	<body>
		<h2> How to get certificate for HTTPS?</h2>
		<p class="ridge">
			<code>
				# run this in your terminal<br>
				cat &gt ca.pem &lt&ltEOF
                <br>
				{{.CACert}}
				EOF<br>
				cat &gt client.pem &lt&ltEOF
                <br>
				{{.ClientCert}}
				EOF<br>
				<br>
				cat &gt client-key.pem &lt&ltEOF 
                <br>
				{{.ClientKey}}
				EOF 
            </code>
		</p>
		you will see ca.pem, client.pem and client-key.pem in your current fold.<br><br>
		
		<h2> How to get certificate for HTTPS?</h2>
		<p class="ridge">
			<code>
				# run in windows cmd.exe. Set password as you like if necessary<br>
				openssl pkcs12 -export -clcerts -in client.pem -inkey client-key.pem -out client.p12<br><br>
			</code>
		</p>
		
		<h2> How to use it?</h2>
		<p class="ridge">
			<code>
				curl -k --cacert ca.pem  --cert client.pem --key client-key.pem https://$SERVER:{{.Port}}/...<br>
			</code>
		</p>
		<br>
		<p class="double">
			<B>
				Warning: take care of your certificate file, which is important to ensure cluster security !<br>
			</B>
		</p>
	</body>
</html>
`
