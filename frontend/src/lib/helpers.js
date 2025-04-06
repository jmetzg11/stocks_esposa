export async function startSimulation(payLoad) {
	console.log(payLoad);
	try {
		const url = `${import.meta.env.VITE_API_URL}/start_simulation`;
		const response = await fetch(url, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(payLoad)
		});
		if (!response.ok) {
			console.error('Error starting simulation', response.statusText);
			return false;
		}
		const data = await response.json();
		console.log(data);
		return true;
	} catch (error) {
		console.error('Error starting simulation', error);
		return false;
	}
}

export function defaultValues() {
	return {
		marketCap: null,
		negativeTrend: null,
		OnePercentAmount: null,
		TenPercentAmount: null,
		lastBountLimit: null,
		portfolioProportion: null
	};
}
export const formFields = [
	{ value: 'marketCap', label: 'Min Market Cap ($ Billion)' },
	{ value: 'negativeTrend', label: 'Negative Trend limit (weeks)' },
	{ value: 'onePercentBuy', label: '1% Buy Amount' },
	{ value: 'tenPercentBuy', label: '10% Buy Amount' },
	{ value: 'lastBuyLimit', label: 'Since last investment (days)' },
	{ value: 'portfolioProportion', label: 'Proportion of Portfolio (basis points)' }
];
