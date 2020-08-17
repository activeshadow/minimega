(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-5b06f226"],{"6f55":function(t,e,s){"use strict";s.r(e);var i=function(){var t=this,e=t.$createElement,s=t._self._c||e;return s("div",{staticClass:"content"},[s("b-modal",{attrs:{active:t.expModal.active,"has-modal-card":""},on:{"update:active":function(e){return t.$set(t.expModal,"active",e)}}},[s("div",{staticClass:"modal-card",staticStyle:{width:"25em"}},[s("header",{staticClass:"modal-card-head"},[s("p",{staticClass:"modal-card-title"},[t._v("VM "+t._s(t.expModal.vm.name?t.expModal.vm.name:"unknown"))])]),s("section",{staticClass:"modal-card-body"},[s("p",[t._v("Host: "+t._s(t.expModal.vm.host))]),s("p",[t._v("IPv4: "+t._s(t._f("stringify")(t.expModal.vm.ipv4)))]),s("p",[t._v("CPU(s): "+t._s(t.expModal.vm.cpus))]),s("p",[t._v("Memory: "+t._s(t._f("ram")(t.expModal.vm.ram)))]),s("p",[t._v("Disk: "+t._s(t.expModal.vm.disk))]),s("p",[t._v("Uptime: "+t._s(t._f("uptime")(t.expModal.vm.uptime)))]),s("p",[t._v("Network(s): "+t._s(t._f("lowercase")(t._f("stringify")(t.expModal.vm.networks))))]),s("p",[t._v("Taps: "+t._s(t._f("lowercase")(t._f("stringify")(t.expModal.vm.taps))))])]),s("footer",{staticClass:"modal-card-foot"})])]),s("hr"),s("b-field",{attrs:{position:"is-left"}},[s("p",{staticClass:"control"}),s("h3",[t._v("Experiment: "+t._s(this.$route.params.id))]),s("p")]),t.experimentUser()||t.experimentViewer()?s("b-field",{attrs:{position:"is-right"}},[s("b-autocomplete",{attrs:{placeholder:"Find a VM",icon:"search",data:t.filteredData},on:{select:function(e){return t.filtered=e}},model:{value:t.searchName,callback:function(e){t.searchName=e},expression:"searchName"}},[s("template",{slot:"empty"},[t._v("No results found")])],2),s("p",{staticClass:"control"},[s("button",{staticClass:"button",staticStyle:{color:"#686868"},on:{click:function(e){t.searchName=""}}},[s("b-icon",{attrs:{icon:"window-close"}})],1)]),t._v("\n       \n    "),s("p",{staticClass:"control buttons"},[t.adminUser()?s("b-button",{staticClass:"button is-success",attrs:{slot:"trigger","icon-right":"play"},on:{click:t.start},slot:"trigger"}):t._e()],1),t._v("\n       \n    "),s("p",{staticClass:"control"},[s("b-tooltip",{attrs:{label:"menu for scheduling hosts to the experiment",type:"is-light",multilined:""}},[s("b-dropdown",{staticClass:"is-right",attrs:{"aria-role":"list"},model:{value:t.algorithm,callback:function(e){t.algorithm=e},expression:"algorithm"}},[s("button",{staticClass:"button is-light",attrs:{slot:"trigger"},slot:"trigger"},[s("b-icon",{attrs:{icon:"bars"}})],1),t._l(t.schedules,(function(e,i){return s("b-dropdown-item",{key:i,attrs:{value:e},on:{click:t.updateSchedule}},[s("font",{attrs:{color:"#202020"}},[t._v(t._s(e))])],1)}))],2)],1)],1)],1):t._e(),s("div",[s("b-tabs",{on:{change:t.updateFiles}},[s("b-tab-item",{attrs:{label:"Table"}},[s("b-table",{key:t.table.key,attrs:{data:t.vms,paginated:t.table.isPaginated&&t.paginationNeeded,"per-page":t.table.perPage,"current-page":t.table.currentPage,"pagination-simple":t.table.isPaginationSimple,"pagination-size":t.table.paginationSize,"default-sort-direction":t.table.defaultSortDirection,"default-sort":"name"},on:{"update:currentPage":function(e){return t.$set(t.table,"currentPage",e)},"update:current-page":function(e){return t.$set(t.table,"currentPage",e)}},scopedSlots:t._u([{key:"default",fn:function(e){return[s("b-table-column",{attrs:{field:"name",label:"VM",sortable:""}},[t.adminUser()?[s("b-tooltip",{attrs:{label:"get info on the vm",type:"is-dark"}},[s("div",{staticClass:"field"},[s("div",{on:{click:function(s){t.expModal.active=!0,t.expModal.vm=e.row}}},[t._v("\n                        "+t._s(e.row.name)+"\n                      ")])])])]:[t._v("\n                  "+t._s(e.row.name)+"\n                ")]],2),s("b-table-column",{attrs:{field:"host",label:"Host",width:"200",sortable:""}},[t.adminUser()?[s("b-tooltip",{attrs:{label:"assign the vm to a specific host",type:"is-dark"}},[s("b-field",[s("b-select",{attrs:{value:e.row.host,expanded:""},on:{input:function(s){return t.assignHost(e.row.name,s)}}},t._l(t.hosts,(function(e,i){return s("option",{key:i,domProps:{value:e}},[t._v("\n                          "+t._s(e)+"\n                        ")])})),0),s("p",{staticClass:"control"},[s("button",{staticClass:"button",on:{click:function(s){return t.unassignHost(e.row.name)}}},[s("b-icon",{attrs:{icon:"window-close"}})],1)])],1)],1)]:[t._v("\n                  "+t._s(e.row.host)+"\n                ")]],2),s("b-table-column",{attrs:{field:"ipv4",label:"IPv4"}},t._l(e.row.ipv4,(function(e){return s("div",[t._v("\n                  "+t._s(e)+"\n                ")])})),0),s("b-table-column",{attrs:{field:"cpus",label:"CPUs",sortable:"",centered:""}},[t.adminUser()?[s("b-tooltip",{attrs:{label:"menu for assigning vm(s) cpus",type:"is-dark"}},[s("b-select",{attrs:{value:e.row.cpus,expanded:""},on:{input:function(s){return t.assignCpu(e.row.name,s)}}},[s("option",{attrs:{value:"1"}},[t._v("1")]),s("option",{attrs:{value:"2"}},[t._v("2")]),s("option",{attrs:{value:"3"}},[t._v("3")]),s("option",{attrs:{value:"4"}},[t._v("4")])])],1)]:[t._v("\n                  "+t._s(e.row.cpus)+"\n                ")]],2),s("b-table-column",{attrs:{field:"ram",label:"Memory",sortable:"",centered:""}},[t.adminUser()?[s("b-tooltip",{attrs:{label:"menu for assigning vm(s) memory",type:"is-dark"}},[s("b-select",{attrs:{value:e.row.ram,expanded:""},on:{input:function(s){return t.assignRam(e.row.name,s)}}},[s("option",{attrs:{value:"512"}},[t._v("512 MB")]),s("option",{attrs:{value:"1024"}},[t._v("1 GB")]),s("option",{attrs:{value:"2048"}},[t._v("2 GB")]),s("option",{attrs:{value:"3072"}},[t._v("3 GB")]),s("option",{attrs:{value:"4096"}},[t._v("4 GB")]),s("option",{attrs:{value:"8192"}},[t._v("8 GB")]),s("option",{attrs:{value:"12288"}},[t._v("12 GB")]),s("option",{attrs:{value:"16384"}},[t._v("16 GB")])])],1)]:[t._v("\n                  "+t._s(e.row.ram)+"\n                ")]],2),s("b-table-column",{attrs:{field:"disk",label:"Disk"}},[t.adminUser()?[s("b-tooltip",{attrs:{label:"menu for assigning vm(s) disk",type:"is-dark"}},[s("b-select",{attrs:{value:e.row.disk},on:{input:function(s){return t.assignDisk(e.row.name,s)}}},t._l(t.disks,(function(e,i){return s("option",{key:i,domProps:{value:e}},[t._v("\n                          "+t._s(e)+"\n                      ")])})),0)],1)]:[t._v("\n                  "+t._s(e.row.disk)+"\n                ")]],2),t.experimentUser()?s("b-table-column",{attrs:{label:"Boot",centered:""}},[s("b-tooltip",{attrs:{label:"control whether or not VM boots",type:"is-dark"}},[s("div",{on:{click:function(s){return t.updateDnb(e.row.name,!e.row.dnb)}}},[s("font-awesome-icon",{class:t.bootDecorator(e.row.dnb),attrs:{icon:"bolt"}})],1)])],1):t._e()]}}])},[s("template",{slot:"empty"},[s("section",{staticClass:"section"},[s("div",{staticClass:"content has-text-white has-text-centered"},[t._v("\n                  Your search turned up empty!\n                ")])])])],2),s("br"),t.paginationNeeded?s("b-field",{attrs:{grouped:"",position:"is-right"}},[s("div",{staticClass:"control is-flex"},[s("b-switch",{attrs:{size:"is-small",type:"is-light"},model:{value:t.table.isPaginated,callback:function(e){t.$set(t.table,"isPaginated",e)},expression:"table.isPaginated"}},[t._v("Pagenate")])],1)]):t._e()],1),s("b-tab-item",{attrs:{label:"Files"}},[t.files&&!t.files.length?[s("section",{staticClass:"hero is-light is-bold is-large"},[s("div",{staticClass:"hero-body"},[s("div",{staticClass:"container",staticStyle:{"text-align":"center"}},[s("h1",{staticClass:"title"},[t._v("\n                  There are no files available.\n                ")])])])])]:[s("ul",{staticClass:"fa-ul",staticStyle:{"list-style":"none"}},t._l(t.files,(function(e,i){return s("li",{key:i},[s("font-awesome-icon",{staticClass:"fa-li",attrs:{icon:"file-download"}}),s("a",{attrs:{href:"/api/v1/experiments/"+t.experiment.name+"/files/"+e+"?token="+t.$store.state.token,target:"_blank"}},[t._v("\n                "+t._s(e)+"\n              ")])],1)})),0)]],2)],1)],1),s("b-loading",{attrs:{"is-full-page":!0,active:t.isWaiting,"can-cancel":!1},on:{"update:active":function(e){t.isWaiting=e}}})],1)},n=[],a=(s("a481"),s("75fc")),o=(s("28a5"),s("ac6a"),s("6762"),s("2fdb"),s("6b54"),s("7f7f"),s("4917"),s("3b2b"),s("96cf"),s("3b8d")),r={beforeDestroy:function(){this.$options.sockets.onmessage=null},created:function(){var t=Object(o["a"])(regeneratorRuntime.mark((function t(){return regeneratorRuntime.wrap((function(t){while(1)switch(t.prev=t.next){case 0:this.$options.sockets.onmessage=this.handler,this.updateExperiment(),this.adminUser()&&(this.updateHosts(),this.updateDisks());case 3:case"end":return t.stop()}}),t,this)})));function e(){return t.apply(this,arguments)}return e}(),computed:{vms:function(){var t=this.experiment.vms,e=new RegExp(this.searchName,"i"),s=[];for(var i in t){var n=t[i];n.name.match(e)&&s.push(n)}return s},filteredData:function(){var t=this,e=this.vms.map((function(t){return t.name}));return e.filter((function(e){return e.toString().toLowerCase().indexOf(t.searchName.toLowerCase())>=0}))},paginationNeeded:function(){return!(this.vms.length<=10)}},methods:{adminUser:function(){return["Global Admin","Experiment Admin"].includes(this.$store.getters.role)},experimentUser:function(){return["Global Admin","Experiment Admin","Experiment User"].includes(this.$store.getters.role)},experimentViewer:function(){return["Experiment Viewer"].includes(this.$store.getters.role)},bootDecorator:function(t){return t?"":"boot"},handler:function(t){var e=this;t.data.split(/\r?\n/).forEach((function(t){var s=JSON.parse(t);e.handle(s)}))},handle:function(t){switch(t.resource.type){case"experiment":if("schedule"!=t.resource.action)return;for(var e=this.experiment.vms,s=0;s<t.result.schedule.length;s++)for(var i=0;s<e.length;i++)if(e[i].name==t.result.schedule[s].vm){e[i].host=t.result.schedule[s].host;break}this.experiment.vms=Object(a["a"])(e),this.$buefy.toast.open({message:"The VMs for this experiment have been scheduled.",type:"is-success"});break;case"experiment/vm":if("update"!=t.resource.action)return;for(var n=this.experiment.vms,o=0;o<n.length;o++)if(n[o].name==t.result.name){n[o]=t.result;break}this.experiment.vms=Object(a["a"])(n),this.$buefy.toast.open({message:"The VM "+t.result.name+" has been successfully updated.",type:"is-success"});break}},updateExperiment:function(){var t=this;this.$http.get("experiments/"+this.$route.params.id).then((function(e){e.json().then((function(e){t.experiment=e,t.isWaiting=!1}))}),(function(e){t.isWaiting=!1,t.$buefy.toast.open({message:"Getting the experiments failed.",type:"is-danger",duration:4e3})}))},updateHosts:function(){var t=this;this.$http.get("hosts").then((function(e){e.json().then((function(e){if(0==e.hosts.length)t.isWaiting=!0;else{for(var s=0;s<e.hosts.length;s++)e.hosts[s].schedulable&&t.hosts.push(e.hosts[s].name);t.isWaiting=!1}}))}),(function(e){t.isWaiting=!1,t.$buefy.toast.open({message:"Getting the hosts failed.",type:"is-danger",duration:4e3})}))},updateDisks:function(){var t=this;this.$http.get("disks").then((function(e){e.json().then((function(e){if(0==e.disks.length)t.isWaiting=!0;else{for(var s=0;s<e.disks.length;s++)t.disks.push(e.disks[s]);t.isWaiting=!1}}))}),(function(e){t.isWaiting=!1,t.$buefy.toast.open({message:"Getting the disks failed.",type:"is-danger",duration:4e3})}))},updateFiles:function(){var t=this;this.$http.get("experiments/"+this.$route.params.id+"/files").then((function(e){e.json().then((function(e){for(var s=0;s<e.files.length;s++)t.files.push(e.files[s]);t.isWaiting=!1}))}),(function(e){t.isWaiting=!1,t.$buefy.toast.open({message:"Getting the files failed.",type:"is-danger",duration:4e3})}))},start:function(){var t=this;this.$buefy.dialog.confirm({title:"Start the Experiment",message:"This will start the "+this.$route.params.id+" experiment.",cancelText:"Cancel",confirmText:"Start",type:"is-success",hasIcon:!0,onConfirm:function(){t.isWaiting=!0,t.$http.post("experiments/"+t.$route.params.id+"/start").then((function(e){console.log("the "+t.$route.params.id+" experiment was started."),t.$router.replace("/experiments/")}),(function(e){t.$buefy.toast.open({message:"Starting experiment "+t.$route.params.id+" failed with "+e.status+" status.",type:"is-danger",duration:4e3}),t.isWaiting=!1}))}})},assignHost:function(t,e){var s=this;this.$buefy.dialog.confirm({title:"Assign a Host",message:"This will assign the "+t+" VM to the "+e+" host.",cancelText:"Cancel",confirmText:"Assign Host",type:"is-success",hasIcon:!0,onConfirm:function(){s.isWaiting=!0;var i={host:e};s.$http.patch("experiments/"+s.$route.params.id+"/vms/"+t,i).then((function(t){for(var e=s.experiment.vms,i=0;i<e.length;i++)if(e[i].name==t.body.name){e[i]=t.body;break}s.experiment.vms=Object(a["a"])(e),s.isWaiting=!1}),(function(i){s.$buefy.toast.open({message:"Assigning the "+t+" VM to the "+e+" host failed with "+i.status+" status.",type:"is-danger",duration:4e3}),s.isWaiting=!1}))},onCancel:function(){s.tableKey+=1}})},unassignHost:function(t){var e=this;this.$buefy.dialog.confirm({title:"Unassign a Host",message:"This will cancel the host assignment for "+t+" VM.",cancelText:"Cancel",confirmText:"Unassign Host",type:"is-success",hasIcon:!0,onConfirm:function(){e.isWaiting=!0;var s={host:""};e.$http.patch("experiments/"+e.$route.params.id+"/vms/"+t,s).then((function(t){for(var s=e.experiment.vms,i=0;i<s.length;i++)if(s[i].name==t.body.name){s[i]=t.body;break}e.experiment.vms=Object(a["a"])(s),e.isWaiting=!1}),(function(s){e.$buefy.toast.open({message:"Canceling the "+host+" assignment for the "+t+" VM failed with "+s.status+" status.",type:"is-danger",duration:4e3}),e.isWaiting=!1}))}})},assignCpu:function(t,e){var s=this;this.$buefy.dialog.confirm({title:"Assign CPUs",message:"This will assign "+e+" cpu(s) to the "+t+" VM.",cancelText:"Cancel",confirmText:"Assign CPUs",type:"is-success",hasIcon:!0,onConfirm:function(){s.isWaiting=!0;var i={cpus:e};s.$http.patch("experiments/"+s.$route.params.id+"/vms/"+t,i).then((function(t){for(var e=s.experiment.vms,i=0;i<e.length;i++)if(e[i].name==t.body.name){e[i]=t.body;break}s.experiment.vms=Object(a["a"])(e),s.isWaiting=!1}),(function(i){s.$buefy.toast.open({message:"Assigning "+e+" cpu(s) to the "+t+" VM failed with "+i.status+" status.",type:"is-danger",duration:4e3}),s.isWaiting=!1}))},onCancel:function(){s.tableKey+=1}})},assignRam:function(t,e){var s=this;this.$buefy.dialog.confirm({title:"Assign Memory",message:"This will assign "+e+" of memory to the "+t+" VM.",cancelText:"Cancel",confirmText:"Assign Memory",type:"is-success",hasIcon:!0,onConfirm:function(){s.isWaiting=!0;var i={ram:e};s.$http.patch("experiments/"+s.$route.params.id+"/vms/"+t,i).then((function(t){for(var e=s.experiment.vms,i=0;i<e.length;i++)if(e[i].name==t.body.name){e[i]=t.body;break}s.experiment.vms=Object(a["a"])(e),s.isWaiting=!1}),(function(i){s.$buefy.toast.open({message:"Assigning "+e+" of memory to the "+t+" VM failed with "+i.status+" status.",type:"is-danger",duration:4e3}),s.isWaiting=!1}))},onCancel:function(){s.tableKey+=1}})},assignDisk:function(t,e){var s=this;this.$buefy.dialog.confirm({title:"Assign a Disk Image",message:"This will assign the "+e+" disk image to the "+t+" VM.",cancelText:"Cancel",confirmText:"Assign Disk",type:"is-success",hasIcon:!0,onConfirm:function(){s.isWaiting=!0;var i={disk:e};s.$http.patch("experiments/"+s.$route.params.id+"/vms/"+t,i).then((function(t){for(var e=s.experiment.vms,i=0;i<e.length;i++)if(e[i].name==t.body.name){e[i]=t.body;break}s.experiment.vms=Object(a["a"])(e),s.isWaiting=!1}),(function(i){s.$buefy.toast.open({message:"Assigning the "+e+" to the "+t+" VM failed with "+i.status+" status.",type:"is-danger",duration:4e3}),s.isWaiting=!1}))},onCancel:function(){s.tableKey+=1}})},updateDnb:function(t,e){var s=this;e?this.$buefy.dialog.confirm({title:"Set Do NOT Boot",message:"This will set the "+t+" VM to NOT boot when the experiment starts.",cancelText:"Cancel",confirmText:"Do NOT Boot",type:"is-warning",hasIcon:!0,onConfirm:function(){s.isWaiting=!0;var i={dnb:e};s.$http.patch("experiments/"+s.$route.params.id+"/vms/"+t,i).then((function(t){for(var e=s.experiment.vms,i=0;i<e.length;i++)if(e[i].name==t.body.name){e[i]=t.body;break}s.experiment.vms=Object(a["a"])(e),s.isWaiting=!1}),(function(e){s.$buefy.toast.open({message:"Setting the "+t+" VM to NOT boot when experiment starts failed with "+e.status+" status.",type:"is-danger",duration:4e3}),s.isWaiting=!1}))}}):this.$buefy.dialog.confirm({title:"Set Boot",message:"This will set the "+t+" VM to boot when the experiment starts.",cancelText:"Cancel",confirmText:"Boot",type:"is-success",hasIcon:!0,onConfirm:function(){s.isWaiting=!0;var i={dnb:e};s.$http.patch("experiments/"+s.$route.params.id+"/vms/"+t,i).then((function(t){for(var e=s.experiment.vms,i=0;i<e.length;i++)if(e[i].name==t.body.name){e[i]=t.body;break}s.experiment.vms=Object(a["a"])(e),s.isWaiting=!1}),(function(e){s.$buefy.toast.open({message:"Setting the "+t+" VM to boot when experiment starts failed with "+e.status+" status.",type:"is-danger",duration:4e3}),s.isWaiting=!1}))}})},updateSchedule:function(){var t=this;this.$buefy.dialog.confirm({title:"Assign a Host Schedule",message:"This will schedule host(s) with the "+this.algorithm+" algorithm for the "+this.$route.params.id+" experiment.",cancelText:"Cancel",confirmText:"Assign Schedule",type:"is-success",hasIcon:!0,onConfirm:function(){t.isWaiting=!0,t.$http.post("experiments/"+t.$route.params.id+"/schedule",{algorithm:t.algorithm}).then((function(e){for(var s=t.experiment.vms,i=0;i<s.length;i++)if(s[i].name==e.body.name){s[i]=e.body;break}t.experiment.vms=Object(a["a"])(s),t.isWaiting=!1}),(function(e){t.$buefy.toast.open({message:"Scheduling the host(s) with the "+t.algorithm+" for the "+t.$route.params.id+" experiment failed with "+e.status+" status.",type:"is-danger",duration:4e3}),t.isWaiting=!1}))}})}},data:function(){return{table:{key:0,isPaginated:!0,perPage:10,currentPage:1,isPaginationSimple:!0,paginationSize:"is-small",defaultSortDirection:"asc"},expModal:{active:!1,vm:[]},schedules:["isolate_experiment","round_robin"],experiment:[],files:[],hosts:[],disks:[],searchName:"",filtered:null,algorithm:null,dnb:!1,isWaiting:!0}}},l=r,c=(s("763f"),s("2877")),u=Object(c["a"])(l,i,n,!1,null,"84d9e842",null);e["default"]=u.exports},"763f":function(t,e,s){"use strict";var i=s("8927"),n=s.n(i);n.a},8927:function(t,e,s){}}]);
//# sourceMappingURL=chunk-5b06f226.89fd3532.js.map