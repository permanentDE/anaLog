package frontend

var files = map[string]string{
	"index.html": `
<!DOCTYPE html>
<html lang="de">
	<head>
		<meta charset="utf-8">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<link rel="stylesheet" href="//maxcdn.bootstrapcdn.com/bootstrap/3.3.1/css/bootstrap.min.css">
		<link rel="stylesheet" href="//maxcdn.bootstrapcdn.com/bootstrap/3.3.1/css/bootstrap-theme.min.css">

		<title>anaLog Overview</title>

		<script type="text/javascript" src="//code.jquery.com/jquery-1.11.3.min.js"></script>
		<script type="text/javascript">
			function debounce(fn, delay) {
				var timer = null;
				return function () {
					var context = this, args = arguments;
					clearTimeout(timer);
					timer = setTimeout(function () {
						fn.apply(context, args);
					}, delay);
				};
			}

			function reloadLogOverviewValues() {
				task = $('#filterTask').val()
				runid = $('#filterRunId').val()
				host = $('#filterHost').val()
				state = $('#filterState').val()
				rawRegex = $('#filterRaw').val()
				limit = $('#filterLimit').val()

				data = {	
					"admin-secret": adminsecret,
					"task": task,
					"runId": runid,
					"host": host,
					"state": state,
					"rawRegex": rawRegex
				}

				beginDate = document.getElementById("filterDateBegin").valueAsDate
				beginTime = document.getElementById("filterTimeBegin").valueAsDate
				endDate = document.getElementById("filterDateEnd").valueAsDate
				endTime = document.getElementById("filterTimeEnd").valueAsDate

				if(beginDate) {
					begin = new Date();
					begin.setFullYear(beginDate.getFullYear())
					begin.setMonth(beginDate.getMonth())
					begin.setDate(beginDate.getDate())
					if (beginTime) {
						begin.setHours(beginTime.getHours()-1)
						begin.setMinutes(beginTime.getMinutes())
					} else {
						begin.setHours(0)
						begin.setMinutes(0)
					}
					data.timeRangeGTE = Math.round(begin.getTime()/1000)
				}

				if(endDate) {
					end = new Date();
					end.setFullYear(endDate.getFullYear())
					end.setMonth(endDate.getMonth())
					end.setDate(endDate.getDate())
					if (endTime) {
						end.setHours(endTime.getHours()-1)
						end.setMinutes(endTime.getMinutes())
					}else {
						begin.setHours(23)
						begin.setMinutes(59)
					}
					data.timeRangeLTE = Math.round(end.getTime()/1000)
				}

				$.post("/v1/read/find/"+limit, data, function(data){
					$('.logOverview > .valueContainer').empty();
					$.each(data, function(i,e){
						if(((i % 2) == 0)) {
							loi = logOverviewItemB
						} else {
							loi = logOverviewItemA
						}
						ele = $(loi.replace("||TASK||", e.Task).replace("||HOST||", e.Host).replace("||STATE||", e.State).replace("||TIME||", e.Time))
						ele.attr("data-logpoint", JSON.stringify(e))
						ele.click(function() {
							logpoint = JSON.parse($(this).attr("data-logpoint"))
							$.post("/v1/read/find/0", {"admin-secret": adminsecret, "task": logpoint.Task, "runId": logpoint.RunId}, function(lps) {
								$('#fullView').empty()
								$('#fullView').append(twoCols("Task", logpoint.Task))
								$('#fullView').append(twoCols("RunId", logpoint.RunId))
								$('#fullView').append($('<div class="spacer"></div>'))
								dual = $('<div class="row"></div>')
								$.each(lps, function(index, lp) {
									lpView = $('<div class="col-md-6"></div>')
									lpView.append(twoCols("Host", lp.Host))
									lpView.append(twoCols("Mode", lp.Mode))
									lpView.append(twoCols("State", lp.State))
									lpView.append(twoCols("Time", lp.Time))
									lpView.append(twoCols("Message", lp.Message))
									lpView.append(twoCols("Raw", "<pre>"+lp.Raw+"</pre>"))

									$('#fullView').append(dual)
									dual.append(lpView)

									if(index != 0 && ((index % 2) == 0)) {
										dual = $('<div class="row"></div>')
									}
								})
								btn = $('<button style="width:100%">Close</button>')
								btn.click(function() {$('#fullView').empty()})
								$('#fullView').append(btn)
							}, "json")
						})
						$('.logOverview > .valueContainer').append(ele)
					})
				}, "json")
			}

			function twoCols(key, val) {
				return $(` + "`" + `
					<div class="row">
						<div class="col-md-2">
							` + "`" + `+key+` + "`" + `
						</div>
						<div class="col-md-10">
							` + "`" + `+val+` + "`" + `
						</div>
					</div>
				` + "`" + `)
			}

			var adminsecret = prompt("Admin?", "secret");

			var logOverview = ` + "`" + `
				<div class="container logOverview">
					<div class="row center bold">
						<div class="col-md-3 name cola">
							Task
						</div>
						<div class="col-md-3 name colb">
							Host
						</div>
						<div class="col-md-3 name cola">
							State
						</div>
						<div class="col-md-3 name colb">
							Time
						</div>
					</div>
					<div class="valueContainer">
					</div>
				</div>
			` + "`" + `;

			var logOverviewItemA = ` + "`" + `
				<div class="spacer"></div>
				<div class="row point">
					<div class="col-md-3 value cola">
						||TASK||
					</div>
					<div class="col-md-3 value colb">
						||HOST||
					</div>
					<div class="col-md-3 value cola">
						||STATE||
					</div>
					<div class="col-md-3 value colb">
						||TIME||
					</div>
				</div>
			` + "`" + `;

			var logOverviewItemB = ` + "`" + `
				<div class="spacer"></div>
				<div class="row point">
					<div class="col-md-3 value colb">
						||TASK||
					</div>
					<div class="col-md-3 value cola">
						||HOST||
					</div>
					<div class="col-md-3 value colb">
						||STATE||
					</div>
					<div class="col-md-3 value cola">
						||TIME||
					</div>
				</div>
			` + "`" + `;

			$().ready(function() {
				$('body').append($(logOverview))

				$('#filterForm input').keypress(function() {$(this).change()})
				$('#filterForm input').change(debounce(reloadLogOverviewValues, 300))
				$('#filterForm select').change(reloadLogOverviewValues)
				reloadLogOverviewValues()
				nagiosStatus()
				reloadProblems()
			})

			function nagiosStatus() {
				$.get("/v1/nagios/status", function(data) {
					$('#nagios-status').html(data)
					if (data.indexOf("NAGIOS_OK") < 0) {
						$('#nagios-status').css("background-color", "red")
						btn = $('<button>Reset</button>')
						btn.click(function(){
							$.get("/v1/nagios/reset", {"admin-secret": adminsecret}, function() {
								nagiosStatus()
							})
						})
						$('#nagios-status').append(btn)
					} else {
						$('#nagios-status').css("background-color", "green")
					}
					setTimeout(nagiosStatus, 30000)
				})
			}

			function reloadProblems() {
				$('#problems').empty()

				$.get("/v1/read/problems", {"admin-secret": adminsecret}, function(data) {
					var found = false
					for (var key in data) {
						found = true
						$('#problems').append($(` + "`" + `
							<div class="row">
								<div class="col-md-4">
									` + "`" + `+key+` + "`" + `
								</div>
								<div class="col-md-8 ">
									` + "`" + `+data[key]+` + "`" + `
								</div>
							</div>
						` + "`" + `))
					}
					if (found == true) {
						$('#problems').prepend($(` + "`" + `
							<div class="row">
								<div class="col-md-4 bold">
									Task
								</div>
								<div class="col-md-8 bold">
									Problem
								</div>
							</div>
						` + "`" + `))
					}
				}, "json")
				
				setTimeout(reloadProblems, 30000)

			}
		</script>

		<style type="text/css">
			html {
				position: relative;
				min-height: 100%;
			}

			body {
				/* Margin bottom by footer height */
				margin-bottom: 60px;
			}

			/*body > .container {
				padding: 60px 15px 0;
			}*/

			.container .text-muted {
				margin: 20px 0;
			}

			.footer > .container {
				padding-right: 15px;
				padding-left: 15px;
			}

			.footer {
				position: absolute;
				bottom: 0;
				width: 100%;
				/* Set the fixed height of the footer here */
				height: 60px;
				background-color: #f5f5f5;
			}

			.bold {
				font-weight: 900;
			}

			.center {				
				text-align: center;
			}

			.logOverview {
				background-color: lightgrey;
			}

			.spacer {
				height: 10px;
			}

			.row > .cola {
				background-color: #dfdfdf;
			}

			.row > .colb {
				background-color: #cfcfcf;
			}

			#nagios-status {
				text-align: center;
			}

			#problems {
				background-color: #ffdfdf;
			}

			#filterState {
				height: 26px;
			}

		</style>
	</head>
	<body>
		<div class="spacer"></div>
		<div class="container" id="nagios-status">
			NAGIOS STATUS
		</div>
		<div class="container" id="problems">
		</div>
		<div class="spacer"></div>
		<div class="container" id="fullView">
		</div>
		<div class="spacer"></div>
		<div class="container">
			<form id="filterForm">
				<div class="row">
					<input class="col-md-3" id="filterTask" type="text" placeholder="Task">
					<input class="col-md-3" id="filterHost" type="text" placeholder="Host">
					<select class="col-md-3" id="filterState" name="filterState">
						<option value="" selected>State</option>
						<option value="Started">Started</option>
						<option value="Running">Running</option>
						<option value="OK">OK</option>
						<option value="Failed">Failed</option>
						<option value="CompletedWithError">Completed with Error</option>
						<option value="Unknown">Unknown</option>
					</select>
					<input class="col-md-3" id="filterRaw" type="text" placeholder="Log Regex" >
				</div>
				<div class="row">
					<input class="col-md-2" id="filterLimit" type="text" value="10" placeholder="Limit">			
					<input class="col-md-10" id="filterRunId" type="text" placeholder="RunId">
				</div>
				<div class="row">
					<label class="col-md-2">Datetime Begin</label>
					<div class="col-md-10">
						<input class="col-md-6" type="date" id="filterDateBegin">
						<input class="col-md-6" type="time" id="filterTimeBegin">
					</div>
				</div>
				<div class="row">
					<label class="col-md-2">Datetime End</label>
					<div class="col-md-10">
						<input class="col-md-6" type="date" id="filterDateEnd">
						<input class="col-md-6" type="time" id="filterTimeEnd">
					</div>
				</div>
			</form>
		</div>
		<div class="spacer"></div>

		<footer class="footer">
			<div class="container">
				<p class="text-muted" style="float:right;">2015 &copy; permanent. Wirtschaftsförderung GmbH &amp; Co. KG</p>
			</div>
		</footer>
	</body>
</html>`,
}

func File(name string) (string, bool) {
	str, ok := files[name]
	return str, ok
}
