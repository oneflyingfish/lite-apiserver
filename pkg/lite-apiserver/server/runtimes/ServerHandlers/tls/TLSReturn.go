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

	err = temp.Execute(buffer, info)
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

        <link href="https://cdnjs.cloudflare.com/ajax/libs/prism/1.22.0/themes/prism.min.css" rel="stylesheet" />
	</head>

	<body>
      <div style="display: block; margin: 0 auto; width: 50%;">
      <script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.22.0/prism.min.js"></script>
          <h2> 1. How to get certificate for HTTPS?</h2>
		  <a href="/tls?format=json">view json format</a>
<pre><code class="language-bash"># run this in your terminal
cat &gt ca.pem &lt&ltEOF
{{.CACert}}EOF

cat &gt client.pem &lt&ltEOF
{{.ClientCert}}EOF

cat &gt client-key.pem &lt&ltEOF 
{{.ClientKey}}EOF</code></pre>
		you will see ca.pem, client.pem and client-key.pem in your current fold.<br>
		
          <h2>2. How to get certificate for HTTPS?</h2>
      <pre><code class="language-bash"># run in windows cmd.exe. Set password as you like if necessary
openssl pkcs12 -export -clcerts -in client.pem -inkey client-key.pem -out client.p12</code></pre>

          <h2>3. How to use it?</h2>
          <pre><code class="language-bash">curl -k --cacert ca.pem  --cert client.pem --key client-key.pem https://$SERVER:{{.Port}}/...</code></pre>
          <p>
            Warning: take care of your certificate file, which is important to ensure cluster security !<br>
          </p>
       </div>
	</body>
</html>
`
