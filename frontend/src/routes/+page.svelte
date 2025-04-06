<script>
	import SubmitForm from '$lib/SubmitForm.svelte';
</script>

<SubmitForm />

<style lang="postcss">
	@reference "tailwindcss";
	:global(html) {
		background-color: theme(--color-gray-100);
	}
</style>
