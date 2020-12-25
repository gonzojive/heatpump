package dashboard

var mainScript = content{
	contentType: javascriptContent,
	text: `Plotly.d3.csv("index.csv", function(err, rows){

	function unpack(rows, key) {
	return rows.map(function(row) { return row[key]; });
  }
  
  
  var copTrace = {
	type: "scatter",
	mode: "lines",
	name: 'COP',
	x: unpack(rows, 'Time'),
	y: unpack(rows, 'COP'),
	line: {color: '#17BECF'}
  }
  
  var dataCOP = [copTrace];
  
  var layoutCOP = {
	title: 'COP over time',
	xaxis: {
	  autorange: true,
	  // range: ['2015-02-17', '2017-02-16'],
	  rangeselector: {buttons: [
		  {
			count: 1,
			label: '1m',
			step: 'month',
			stepmode: 'backward'
		  },
		  {
			count: 6,
			label: '6m',
			step: 'month',
			stepmode: 'backward'
		  },
		  {step: 'all'}
		]},
	  // rangeslider: {range: ['2015-02-17', '2017-02-16']},
	  type: 'date'
	},
	yaxis: {
	  autorange: true,
	  // range: [86.8700008333, 138.870004167],
	  type: 'linear',
	  title: {
		text: 'Coefficient of Performance',
	  },
	}
  };


  var layoutTemps = {
	title: 'Significant Temperature values (C) over time',
	xaxis: {
	  autorange: true,
	  // range: ['2015-02-17', '2017-02-16'],
	  rangeselector: {buttons: [
		  {
			count: 1,
			label: '1m',
			step: 'month',
			stepmode: 'backward'
		  },
		  {
			count: 6,
			label: '6m',
			step: 'month',
			stepmode: 'backward'
		  },
		  {step: 'all'}
		]},
	  // rangeslider: {range: ['2015-02-17', '2017-02-16']},
	  type: 'date'
	},
	yaxis: {
	  autorange: true,
	  // range: [86.8700008333, 138.870004167],
	  type: 'linear',
	  title: {
		text: '°C',
	  },
	}
  };

  var tempTraces = [
	{
		type: "scatter",
		mode: "lines",
		name: 'Ambient Temp',
		x: unpack(rows, 'Time'),
		y: unpack(rows, 'Ambient Temp'),
		line: {color: '#7F7F7F'}
	},
	{
		type: "scatter",
		mode: "lines",
		name: 'Inlet Water Temp',
		x: unpack(rows, 'Time'),
		y: unpack(rows, 'Inlet Temp'),
		line: {color: '#17BECF'}
	},
	{
		type: "scatter",
		mode: "lines",
		name: 'Outlet Water Temp',
		x: unpack(rows, 'Time'),
		y: unpack(rows, 'Outlet Temp'),
		line: {color: '#cf1723'}
	},
	{
		type: "scatter",
		mode: "lines",
		name: 'Target Temp',
		x: unpack(rows, 'Time'),
		y: unpack(rows, 'Target Temp'),
		line: {color: '#9117cf'}
	},
  ];

  
  Plotly.newPlot('copDiv', dataCOP, layoutCOP);
  Plotly.newPlot('tempDiv', tempTraces, layoutTemps);
  })
  `,
}
