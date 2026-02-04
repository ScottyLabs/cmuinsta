<script lang="ts">
    import { onMount } from "svelte";

    let message: string = "Loading...";
    let backendStatus: string = "Checking...";

    onMount(async () => {
        try {
            // Check health
            const healthRes = await fetch("/api/health");
            const healthData = await healthRes.json();
            backendStatus = healthData.status === "ok" ? "Online" : "Error";

            // Get data
            const dataRes = await fetch("/api/data");
            const data = await dataRes.json();
            message = data.message;
        } catch (e) {
            console.error(e);
            message = "Failed to connect to backend.";
            backendStatus = "Offline";
        }
    });
</script>

<main
    class="flex min-h-screen flex-col items-center justify-center bg-gray-100 font-sans text-gray-900"
>
    <div class="w-full max-w-md rounded-xl bg-white p-8 shadow-lg">
        <h1 class="mb-2 text-center text-3xl font-bold text-gray-900">
            CMU Insta
        </h1>

        <div class="mb-6 text-center text-sm text-gray-500 space-y-1">
            <p>
                Frontend: <span class="font-medium text-gray-700"
                    >Svelte + Vite + Tailwind</span
                >
            </p>
            <p>
                Backend:
                <span
                    class="font-bold"
                    class:text-green-600={backendStatus === "Online"}
                    class:text-red-600={backendStatus !== "Online"}
                >
                    {backendStatus}
                </span>
            </p>
        </div>

        <div
            class="rounded-lg border border-blue-200 bg-blue-50 p-4 text-blue-900"
        >
            <strong
                class="mb-1 block text-xs font-semibold uppercase tracking-wider text-blue-700"
                >Message from Go</strong
            >
            <p class="text-lg">{message}</p>
        </div>
    </div>
</main>
