function drawAvanceOptimo()
{
	var data = new google.visualization.DataTable();

	data.addColumn('number', 'Iteraciones');
	data.addColumn('number', 'Z');

	data.addRows(mtz_avance_opt);

	var options = {
		height: 600,
		colors: ['#3bafda', 'blue', '#3fc26b'],
		hAxis: {title: 'Iteraciónes'},
		vAxis: {title: 'Z'}
	};

	var grafica = new google.visualization.AreaChart(document.getElementById('chart-avance-optimo'));
	grafica.draw(data, options);
}