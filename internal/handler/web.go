package handler

import "net/http"

const indexPage = `<!doctype html>
<html lang="zh-CN"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>Quasar Weave</title>
<style>body{font:16px system-ui;background:#101722;color:#e9f0f7;margin:0}main{max-width:960px;margin:2rem auto;padding:0 1rem}section{background:#182536;border:1px solid #2b425a;border-radius:12px;padding:1rem;margin:1rem 0}button{background:#62d2a2;border:0;border-radius:6px;padding:.55rem .8rem;cursor:pointer}pre{white-space:pre-wrap;color:#b8d2e8}.grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(180px,1fr));gap:.7rem}.metric{font-size:1.8rem;color:#62d2a2}</style></head>
<body><main><h1>Quasar Weave</h1><p>射电望远镜阵列观测校准值班台</p><section><h2>阵列摘要</h2><div id="metrics" class="grid"></div><button onclick="refresh()">刷新</button></section><section><h2>当前观测</h2><pre id="observations">加载中...</pre></section><section><h2>待处理告警</h2><pre id="alerts">加载中...</pre></section></main>
<script>async function get(path){const r=await fetch(path);return r.json()}async function refresh(){const h=await get('/api/health');document.getElementById('metrics').innerHTML=Object.entries(h).map(([k,v])=>'<div><small>'+k+'</small><div class="metric">'+v+'</div></div>').join('');document.getElementById('observations').textContent=JSON.stringify(await get('/api/observations'),null,2);document.getElementById('alerts').textContent=JSON.stringify(await get('/api/alerts?state=open'),null,2)}refresh()</script></body></html>`

func (rt *Router) handleWeb(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) {
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(indexPage))
}
