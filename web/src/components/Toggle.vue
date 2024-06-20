<script setup lang="ts">
import type { FeatureToggle } from '@/core/types';
import { ref } from 'vue';

const emit = defineEmits<{
    (e: 'change', enabled: boolean): void
}>()

const props = defineProps<{
    toggle: FeatureToggle
}>()

const toggle = ref(props.toggle)

const updateToggle = (ev: Event) => {
  const cb = ev.target as HTMLInputElement

  toggle.value = {
    featureId: toggle.value.featureId,
    description: toggle.value.description,
    enabled: cb.checked,
  };

  emit('change', toggle.value.enabled)
}
</script>

<template>
    <label class="inline-flex items-center cursor-pointer">
        <input type="checkbox" value="" class="sr-only peer" :checked="toggle.enabled" @change="updateToggle">
        <div class="relative w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 dark:peer-focus:ring-blue-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-blue-600"></div>
        <span class="ms-3 text-sm font-medium text-gray-900 dark:text-gray-300">{{ toggle.featureId }}</span>
    </label>
</template>