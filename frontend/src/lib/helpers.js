export async function startSimulation(data) {
	try {
		const url = `${import.meta.env.VITE_API_URL}/start_simulation`;
		const response = await fetch(url, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(data)
		});
		if (!response.ok) {
			console.error('Error starting simulation', response.statusText);
			return false;
		}
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
	{ value: 'marketCap', label: 'Market Cap ($ Billion)' },
	{ value: 'negativeTrend', label: 'Negative Trend limit (weeks)' },
	{ value: 'OnePercentAmount', label: '1% Buy Amount' },
	{ value: 'TenPercentAmount', label: '10% Buy Amount' },
	{ value: 'lastBountLimit', label: 'Since last investment (days)' },
	{ value: 'portfolioProportion', label: 'Proportion of Portfolio (basis points)' }
];
