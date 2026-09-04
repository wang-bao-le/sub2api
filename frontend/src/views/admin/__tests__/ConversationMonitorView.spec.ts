import { readFileSync } from 'node:fs'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'
import { describe, expect, it } from 'vitest'

const source = readFileSync(resolve(dirname(fileURLToPath(import.meta.url)), '../ConversationMonitorView.vue'), 'utf8')

describe('ConversationMonitorView', () => {
  it('uses the standard settings toggle after its saved configuration first loads', () => {
    expect(source).toContain("import Toggle from '@/components/common/Toggle.vue'")
    expect(source).toContain("const configLoaded=ref(false)")
    expect(source).toContain("finally { configLoaded.value=true }")
    expect(source).toContain('<Toggle v-if="configLoaded" v-model="config.enabled" />')
  })
})
