<script lang="ts">
  import { onMount } from 'svelte';

  import { formatDate } from '../utils/formatDate';

  export let datetime: string;

  const units: [Intl.RelativeTimeFormatUnit, number][] = [
    ['day', 86_400_000],
    ['hour', 3_600_000],
    ['minute', 60_000]
  ];

  // Server and browser time zones differ, so localized text is client-only.
  // Until then the ISO date is shown, which renders the same on both.
  let text = datetime.slice(0, 10);
  let title: string | undefined;

  onMount(() => {
    // English to match the rest of the UI copy; the tooltip follows the browser locale.
    const rtf = new Intl.RelativeTimeFormat('en', { numeric: 'auto' });
    const diff = new Date(datetime).getTime() - Date.now();
    const unit = units.find(([, ms]) => Math.abs(diff) >= ms);
    if (unit) {
      text = rtf.format(Math.round(diff / unit[1]), unit[0]);
    } else {
      text = diff < 0 ? 'just now' : 'in a moment';
    }
    title = formatDate(datetime);
  });
</script>

<time {datetime} {title}>{text}</time>
