// Code generated by qtc from "WafBlocked.html.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line WafBlocked.html.qtpl:1
package views

//line WafBlocked.html.qtpl:1
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line WafBlocked.html.qtpl:1
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line WafBlocked.html.qtpl:1
func StreamBuildWafBlockedPageHtml(qw422016 *qt422016.Writer, policyId int, clientIp string, nodeId string) {
//line WafBlocked.html.qtpl:1
	qw422016.N().S(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1"/>
    <meta name="renderer" content="webkit" />
    <meta name="force-rendering" content="webkit" />
    <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1" />
    <link rel="preconnect" href="https://fonts.gstatic.com" />
    <link href="https://fonts.googleapis.com/css2?family=Red+Hat+Display:wght@500&family=Fira+Mono&family=Ubuntu&display=swap" rel="stylesheet"/>
    <title>1020 | Access denied</title>
    <link href="https://fonts.googleapis.com/icon?family=Material+Icons+Outlined" rel="stylesheet"/>
    <style>:root{--background-color:#fff;--font-color:#000;--font-color-lighter:rgb(87,89,88);--font-size-main:3.545rem;--font-size-description:1.245rem;--box-color:#f2f2f2;--working-color:#137333;--working-color-background:#e6f4ea;--error-color-background:#fce8e6;--error-color:#c5221f;--working-with-error-color:#b05a00;--working-with-error-color-background:#fef7e0;--icon-size:48px;}body{margin:2rem 2rem;font-family:"Red Hat Display",Ubuntu,Roboto,Noto Sans SC,sans-serif;color:var(--font-color);background-color:var(--background-color);}header{margin-left:1rem;}header description{font-family:Ubuntu,Roboto,Noto Sans SC,sans-serif;font-size:var(--font-size-description);line-height:var(--font-size-description);color:var(--font-color-lighter);}header main{font-size:var(--font-size-main);line-height:var(--font-size-main);font-family:Fira Mono,Ubuntu,monospace;}code{font-family:Fira Mono,monospace;}none{display:none;}status{margin-top:2.5rem;display:flex;flex-direction:row;flex-wrap:wrap;justify-content:center;align-items:center;}status > card{background-color:var(--box-color);padding:2rem;margin:1rem 1rem;min-height:3rem;border-radius:9px;flex-grow:1;}status > card.green-card{background-color:var(--working-color-background);}status > card.red-card{background-color:var(--error-color-background);}status > card.yellow-card{background-color:var(--working-with-error-color-background);}status > card main{font-size:calc(var(--font-size-description) + 0.1rem);}.green-text{color:var(--working-color);}.red-text{color:var(--error-color);}.yellow-text{color:var(--working-with-error-color);}status-text,reason p{font-family:Ubuntu,Roboto,Noto Sans SC,sans-serif;}reason p{line-height:125%;}icon{font-size:var(--icon-size) !important;}a{text-decoration:none;color:#1967d2;}reason{display:flex;flex-direction:row;flex-wrap:wrap;justify-content:space-between;align-items:baseline;}reason > *{display:block;margin:1rem;flex-grow:1;max-width:40%;}reason main{font-size:calc(var(--font-size-description) + 0.2rem);font-weight:550;}footer{margin:1rem;color:var(--font-color-lighter);font-size:calc(var(--font-size-description) - 0.4rem);}footer > text{font-size:calc(var(--font-size-description) - 0.6rem);}footer > *{display:block;}@media screen and (max-width:480px){body{margin:6rem 2rem;}:root{--font-size-main:3rem;--font-size-description:1.045rem;}reason > *{max-width:100%;}footer{font-size:calc(var(--font-size-description) - 0.2rem);}footer > text{font-size:calc(var(--font-size-description) - 0.4rem);}}@media screen and (min-width:768px){body{margin:8% 10%;}header > *{display:inline-block;margin-left:1%;}}@media (prefers-color-scheme:dark){:root{--font-color:rgba(255,255,255,0.86);--font-color-lighter:rgba(255,255,255,0.4);--background-color:rgb(0,0,0);--box-color:rgb(40 40 40 / 73%);--working-color-background:#07220f;--error-color-background:#270501;--working-with-error-color-background:#392605;}}</style>
    <script defer>document.head.innerHTML += `)
//line WafBlocked.html.qtpl:1
	qw422016.N().S("`")
//line WafBlocked.html.qtpl:1
	qw422016.N().S(`<link href="https://fonts.googleapis.com/css2?family=Noto+Sans+SC&display=swap" rel="stylesheet">`)
//line WafBlocked.html.qtpl:1
	qw422016.N().S("`")
//line WafBlocked.html.qtpl:1
	qw422016.N().S(`;</script>
</head>

<body>
<header>
    <main>1020</main>
    <description> Access denied </description>
</header>

<status>
    <card class="green-card" id="client-status-card">
        <icon class="material-icons-outlined green-text">web_assets</icon>
        <main>Your Client</main>
        <status-text class="green-text">Working</status-text>
    </card>
    <card class="green-card" id="client-status-card">
        <icon class="material-icons-outlined green-text">cloud</icon>
        <main>AckyCDN Edge Network</main>
        <status-text class="green-text">Working</status-text>
    </card>
    <card class="yellow-card" id="client-status-card">
        <icon class="material-icons-outlined yellow-text">dns</icon>
        <main>Web Server</main>
        <status-text class="yellow-text">Unknown</status-text>
    </card>
</status>
<reason>
    <explain>
        <main>What happened?</main>
        <p>
            A client or browser is blocked by a AckyCDN customer’s Firewall
            Rules.
        </p>
    </explain>
    <howtodo>
        <main>What can I do?</main>
        <p>
            Provide the website owner with a screenshot of the 1020 error message
            you received.
        </p>
    </howtodo>
</reason>
</body>

<footer>
    <provider>Performance & security by <a href="https://www.ackystack.com/">AckyCDN</a>.</provider>
    <br/>
    <text>
        Policy ID: &nbsp;<code> `)
//line WafBlocked.html.qtpl:63
	qw422016.N().D(policyId)
//line WafBlocked.html.qtpl:63
	qw422016.N().S(` </code>
        <br/><br/>
        Your IP: &nbsp;<code> `)
//line WafBlocked.html.qtpl:65
	qw422016.E().S(clientIp)
//line WafBlocked.html.qtpl:65
	qw422016.N().S(` </code>
        <br/><br/>
        ENID: &nbsp;<code> `)
//line WafBlocked.html.qtpl:67
	qw422016.E().S(nodeId)
//line WafBlocked.html.qtpl:67
	qw422016.N().S(` </code>
        <br/><br/>
    </text>
</footer>
</html>
`)
//line WafBlocked.html.qtpl:72
}

//line WafBlocked.html.qtpl:72
func WriteBuildWafBlockedPageHtml(qq422016 qtio422016.Writer, policyId int, clientIp string, nodeId string) {
//line WafBlocked.html.qtpl:72
	qw422016 := qt422016.AcquireWriter(qq422016)
//line WafBlocked.html.qtpl:72
	StreamBuildWafBlockedPageHtml(qw422016, policyId, clientIp, nodeId)
//line WafBlocked.html.qtpl:72
	qt422016.ReleaseWriter(qw422016)
//line WafBlocked.html.qtpl:72
}

//line WafBlocked.html.qtpl:72
func BuildWafBlockedPageHtml(policyId int, clientIp string, nodeId string) string {
//line WafBlocked.html.qtpl:72
	qb422016 := qt422016.AcquireByteBuffer()
//line WafBlocked.html.qtpl:72
	WriteBuildWafBlockedPageHtml(qb422016, policyId, clientIp, nodeId)
//line WafBlocked.html.qtpl:72
	qs422016 := string(qb422016.B)
//line WafBlocked.html.qtpl:72
	qt422016.ReleaseByteBuffer(qb422016)
//line WafBlocked.html.qtpl:72
	return qs422016
//line WafBlocked.html.qtpl:72
}
