<script>
	import FormField from './FormField.svelte';
	import { startSimulation, defaultValues, formFields } from './helpers';

	let formData = defaultValues();

	function handleClick() {
		// startSimulation(formData)
		console.log(formData);
	}
</script>

<div class="container mx-auto mt-2 rounded-lg bg-white px-4 py-6 shadow-md">
	<h2 class="mb-6 text-xl font-bold text-gray-800">Simulation Parameters</h2>
	<div class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
		{#each formFields as field}
			<FormField bind:value={formData[field.value]} label={field.label} />
		{/each}
	</div>

	<div class="mt-8 flex justify-end">
		<button
			onclick={handleClick}
			class="rounded-md bg-blue-600 px-6 py-2 font-medium text-white transition-all hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
			>Start Simulation</button
		>
	</div>
</div>
