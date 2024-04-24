import{_ as Z,r as s,o as $,a as d,b as ee,c as N,d as le,e as t,w as a,n as te,u as o,i as u,F as ae,f as C,s as oe,g as F,h as c,j as ne,k as se,t as ie,l as de,E as z}from"./index-Ce00-jc-.js";import{a as h}from"./axios.config-D7HOseWb.js";const ue={class:"card-header"},re={class:"dialog-footer"},ce={class:"dialog-footer"},pe={class:"dialog-footer"},me={__name:"SecretView",setup(fe){let M=s(!1),U=s([]),E=s(1),y=s(10),H=s(0),_=s(""),i=s(!1),v=s({account:"",isEnable:!0}),Y=window.innerHeight,D=s({height:""}),w=s({}),p=s(!1),f=s({}),r=s(!1);$(()=>{D.value.height=Y-108+"px"});const j=()=>{V()},I=()=>{_.value="",V()},T=()=>{V()},V=()=>{h.get("/secret/paging",{params:{pageNo:E.value,pageSize:y.value,searchTxt:_.value}}).then(function(n){let e=n.data.data;U.value=e.rows,H.value=e.total})};V();const K=n=>{h.get("/secret/"+n.account).then(function(e){w.value=e.data.data,p.value=!0})},A=(n,e,g)=>de(g).format("YYYY-MM-DD HH:mm:ss"),P=n=>n===1?"启用":n===2?"禁用":n,L=()=>{const n=v.value;if(!n||!n.account){z.error("缺失参数");return}h.post("/secret",{account:n.account,isEnable:n.isEnable?1:2}).then(function(e){i.value=!1,V(),z.success(e.data.msg)})},R=n=>{f=s({id:n.id,account:n.account,isEnable:n.isEnable===1}),r.value=!0},q=()=>{let n=f.value;h.put("/secret/enable",{id:n.id,account:n.account,isEnable:n.isEnable?1:2}).then(function(e){r.value=!1,V(),z.success(e.data.msg)})};return(n,e)=>{const g=d("el-input"),m=d("el-button"),b=d("el-table-column"),G=d("el-tag"),J=d("el-table"),O=d("el-pagination"),Q=d("el-card"),x=d("el-form-item"),B=d("el-switch"),k=d("el-form"),S=d("el-dialog"),W=ee("loading");return N(),le(ae,null,[t(Q,{class:"box-card",style:te(o(D))},{header:a(()=>[C("div",ue,[t(g,{modelValue:o(_),"onUpdate:modelValue":e[0]||(e[0]=l=>u(_)?_.value=l:_=l),placeholder:"搜索账号 | 启用 | 禁用",class:"search",size:"large",style:{width:"260px"},"suffix-icon":o(oe),clearable:"",onKeydown:[F(j,["enter"]),F(I,["esc"])]},null,8,["modelValue","suffix-icon"]),t(m,{color:"#337ecc",plain:"",style:{"margin-left":"20px"},onClick:e[1]||(e[1]=l=>u(i)?i.value=!0:i=!0)},{default:a(()=>[c(" 新增账号密钥 ")]),_:1})])]),default:a(()=>[ne((N(),se(J,{data:o(U),"element-loading-text":"Loading...",style:{width:"100%"},stripe:""},{default:a(()=>[t(b,{prop:"id",fixed:"",label:"ID","min-width":"30"}),t(b,{prop:"account",label:"账号"}),t(b,{prop:"isEnable",label:"是否启用"},{default:a(l=>[t(G,{type:l.row.isEnable===1?"success":"danger"},{default:a(()=>[c(ie(P(l.row.isEnable)),1)]),_:2},1032,["type"])]),_:1}),t(b,{prop:"dataCheck",label:"数据校验",width:"300","show-overflow-tooltip":!0}),t(b,{prop:"createTime",formatter:A,label:"创建时间",width:"170"}),t(b,{prop:"updateTime",formatter:A,label:"更新时间",width:"170"}),t(b,{label:"操作",width:"200"},{default:a(l=>[t(m,{plain:"",type:"primary",size:"small",onClick:X=>K(l.row)},{default:a(()=>[c(" 查看密钥 ")]),_:2},1032,["onClick"]),t(m,{type:"warning",plain:"",size:"small",onClick:X=>R(l.row)},{default:a(()=>[c(" 启用/禁用 ")]),_:2},1032,["onClick"])]),_:1})]),_:1},8,["data"])),[[W,o(M)]]),t(O,{style:{"margin-top":"10px",float:"right"},"current-page":o(E),"onUpdate:currentPage":e[2]||(e[2]=l=>u(E)?E.value=l:E=l),"page-size":o(y),"onUpdate:pageSize":e[3]||(e[3]=l=>u(y)?y.value=l:y=l),"page-sizes":[5,10,30,100],layout:"total,prev,pager,next,sizes",total:o(H),onSizeChange:T,onCurrentChange:T},null,8,["current-page","page-size","total"])]),_:1},8,["style"]),t(S,{modelValue:o(i),"onUpdate:modelValue":e[7]||(e[7]=l=>u(i)?i.value=l:i=l),title:"新增账号密钥",width:"600px",center:"",draggable:""},{footer:a(()=>[C("span",re,[t(m,{onClick:e[6]||(e[6]=l=>u(i)?i.value=!1:i=!1)},{default:a(()=>[c("取消")]),_:1}),t(m,{color:"#337ecc",plain:"",onClick:L},{default:a(()=>[c(" 确定 ")]),_:1})])]),default:a(()=>[t(k,{model:o(v),"label-width":"100px"},{default:a(()=>[t(x,{label:"账号"},{default:a(()=>[t(g,{modelValue:o(v).account,"onUpdate:modelValue":e[4]||(e[4]=l=>o(v).account=l),clearable:""},null,8,["modelValue"])]),_:1}),t(x,{label:"是否启用"},{default:a(()=>[t(B,{modelValue:o(v).isEnable,"onUpdate:modelValue":e[5]||(e[5]=l=>o(v).isEnable=l),"inline-prompt":"","active-text":"启用","inactive-text":"禁用"},null,8,["modelValue"])]),_:1})]),_:1},8,["model"])]),_:1},8,["modelValue"]),t(S,{modelValue:o(p),"onUpdate:modelValue":e[11]||(e[11]=l=>u(p)?p.value=l:p=l),title:"查看账号密钥",width:"600px",center:"",draggable:""},{footer:a(()=>[C("span",ce,[t(m,{onClick:e[10]||(e[10]=l=>u(p)?p.value=!1:p=!1)},{default:a(()=>[c("关闭")]),_:1})])]),default:a(()=>[t(k,{model:o(w),"label-width":"100px"},{default:a(()=>[t(x,{label:"账号"},{default:a(()=>[t(g,{modelValue:o(w).account,"onUpdate:modelValue":e[8]||(e[8]=l=>o(w).account=l),clearable:""},null,8,["modelValue"])]),_:1}),t(x,{label:"密钥"},{default:a(()=>[t(g,{modelValue:o(w).secret,"onUpdate:modelValue":e[9]||(e[9]=l=>o(w).secret=l),type:"textarea"},null,8,["modelValue"])]),_:1})]),_:1},8,["model"])]),_:1},8,["modelValue"]),t(S,{modelValue:o(r),"onUpdate:modelValue":e[15]||(e[15]=l=>u(r)?r.value=l:r=l),title:"启用/禁用",width:"600px",center:"",draggable:""},{footer:a(()=>[C("span",pe,[t(m,{onClick:e[14]||(e[14]=l=>u(r)?r.value=!1:r=!1)},{default:a(()=>[c("取消")]),_:1}),t(m,{type:"warning",plain:"",onClick:q},{default:a(()=>[c(" 确定修改 ")]),_:1})])]),default:a(()=>[t(k,{model:o(f),"label-width":"100px"},{default:a(()=>[t(x,{label:"账号"},{default:a(()=>[t(g,{modelValue:o(f).account,"onUpdate:modelValue":e[12]||(e[12]=l=>o(f).account=l),disabled:"",clearable:""},null,8,["modelValue"])]),_:1}),t(x,{label:"是否启用"},{default:a(()=>[t(B,{modelValue:o(f).isEnable,"onUpdate:modelValue":e[13]||(e[13]=l=>o(f).isEnable=l),"inline-prompt":"","active-text":"启用","inactive-text":"禁用"},null,8,["modelValue"])]),_:1})]),_:1},8,["model"])]),_:1},8,["modelValue"])],64)}}},_e=Z(me,[["__scopeId","data-v-f825bca4"]]);export{_e as default};
