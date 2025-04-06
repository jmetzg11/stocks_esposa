<script>
	import FormField from './FormField.svelte';
	import { startSimulation, defaultValues, formFields } from './helpers';

	let formData = defaultValues();

	function handleClick() {
		startSimulation(formData);
		console.log(formData);
	}
</script>

<div class="container mx-auto mt-2 rounded-lg bg-white p-3 shadow-md">
	<div class="items-centerpy-2 flex">
		<h2 class="mr-4 text-xl font-bold text-gray-800">Simulation Parameters</h2>
		<button
			onclick={handleClick}
			class="cursor-pointer rounded-md bg-blue-600 px-2 font-medium text-white transition-all hover:bg-blue-700"
			>Start Simulation</button
		>
	</div>
	<div class="mt-4 grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
		{#each formFields as field}
			<FormField bind:value={formData[field.value]} label={field.label} />
		{/each}
	</div>
</div>
