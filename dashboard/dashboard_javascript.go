package dashboard

var mainScript = content{
	contentType: javascriptContent,
	text: `Plotly.d3.csv("index.csv", function(err, rows){

	function unpack(rows, key) {
	return rows.map(function(row) { return row[key]; });
  }

  const sharedRangeSelector = {
	buttons: [
		{
			count: 1,
			label: '1h',
			step: 'hour',
			stepmode: 'backward'
		},
		{
			count: 24,
			label: '24h',
			step: 'hour',
			stepmode: 'backward'
		},
		{step: 'all'}
	]
  };
  
  
  var copTrace = {
	type: "scatter",
	mode: "lines",
	name: 'COP',
	x: unpack(rows, 'Time'),
	y: unpack(rows, 'COP'),
	line: {color: '#17BECF'}
  }
  
  var dataCOP = [copTrace, {
	type: "scatter",
	mode: "lines",
	name: 'Power Consumption (kW)',
	x: unpack(rows, 'Time'),
	y: unpack(rows, 'Approx Power'),
	line: {color: '#7F7F7F'},
	yaxis: 'y2',
  }];
  
  var layoutCOP = {
	title: 'COP over time',
	xaxis: {
	  autorange: true,
	  // range: ['2015-02-17', '2017-02-16'],
	  rangeselector: sharedRangeSelector,
	  // rangeslider: {range: ['2015-02-17', '2017-02-16']},
	  type: 'date'
	},
	yaxis: {
	  autorange: false,
	  range: [0, 5],
	  // range: [86.8700008333, 138.870004167],
	  type: 'linear',
	  title: {
		text: 'Coefficient of Performance',
	  },
	},
	yaxis2: {
	  autorange: false,
	  range: [0, 4.5],
	  type: 'linear',
	  title: {
		text: 'Power Consumption (kW)',
	  },
	  side: 'right',
	  overlaying: 'y',
	  ticksuffix: ' kW',
	}
  };


  var layoutTemps = {
	title: 'Significant Temperature values (C) over time',
	xaxis: {
	  autorange: true,
	  // range: ['2015-02-17', '2017-02-16'],
	  rangeselector: sharedRangeSelector,
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
	  ticksuffix: ' °C',
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

  var layoutPump = {
	title: 'Flow rate and pump speed over time',
	xaxis: {
	  autorange: true,
	  // range: ['2015-02-17', '2017-02-16'],
	  rangeselector: sharedRangeSelector,
	  // rangeslider: {range: ['2015-02-17', '2017-02-16']},
	  type: 'date'
	},
	yaxis: {
	  autorange: true,
	  // range: [86.8700008333, 138.870004167],
	  type: 'linear',
	  title: {
		text: 'liters/minute',
	  },
	  side: 'left',
	  ticksuffix: ' l/m',
	},
	yaxis2: {
	  autorange: true,
	  // range: [-1, 11],
	  type: 'linear',
	  title: {
		text: 'Pump Speed 1-10',
	  },
	  side: 'right',
	  overlaying: 'y',
	}
  };

  var pumpTraces = [
	{
		type: "scatter",
		mode: "lines",
		name: 'Flow Rate',
		x: unpack(rows, 'Time'),
		y: unpack(rows, 'Flow Rate'),
		line: {color: '#17BECF'},
		yaxis: 'y',
	},
	{
		type: "scatter",
		mode: "lines",
		name: 'Recommended Flow Rate@full power',
		x: unpack(rows, 'Time'),
		y: unpack(rows, 'Flow Rate').map((x) => 20),
		line: {color: '#000000'},
		yaxis: 'y',
	},
	{
		type: "scatter",
		mode: "lines",
		name: 'Pump Speed',
		x: unpack(rows, 'Time'),
		y: unpack(rows, 'Pump Speed'),
		line: {color: '#7F7F7F'},
		yaxis: 'y2',
	},
  ];

  var layoutPowerQuality = {
	title: 'Voltage over time',
	xaxis: {
	  autorange: true,
	  rangeselector: sharedRangeSelector,
	  type: 'date'
	},
	yaxis: {
	  autorange: true,
	  type: 'linear',
	  title: {
		text: 'Line Voltage (V)',
	  },
	  side: 'left',
	  ticksuffix: ' V',
	}
  };

  var powerQualityTraces = [
	{
		type: "scatter",
		mode: "lines",
		name: 'Voltage',
		x: unpack(rows, 'Time'),
		y: unpack(rows, 'Voltage'),
		line: {color: '#17BECF'},
		yaxis: 'y',
	}
  ];

  
  Plotly.newPlot('copDiv', dataCOP, layoutCOP);
  Plotly.newPlot('tempDiv', tempTraces, layoutTemps);
  Plotly.newPlot('pumpDiv', pumpTraces, layoutPump);
  Plotly.newPlot('powerQualityDiv', powerQualityTraces, layoutPowerQuality);
  })
  `,
}
