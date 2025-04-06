<script>
	let { value = $bindable(), ...props } = $props();
</script>

<div class="flex flex-col space-y-2">
	<label class="text-sm font-medium text-gray-700">{props.label}</label>
	<input
		type="number"
		bind:value
		min={props.label.slice(0, 3) === 'Min' ? 1 : undefined}
		class="rounded-md border border-gray-300 px-4 py-2 focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500/50"
	/>
</div>
