import { expect, test } from '@playwright/test'

test('landing page shows the hero and navigation', async ({ page }) => {
  await page.goto('/')
  await expect(page.getByRole('heading', { level: 1 })).toHaveText('Forge Your Hero')
  await expect(page.getByRole('navigation', { name: 'Main navigation' })).toBeVisible()
  await expect(page.getByRole('link', { name: 'Content' })).toBeVisible()
})

test('wizard flow: create a character end to end', async ({ page }) => {
  await page.goto('/#/builder')

  // Class step
  await expect(page.getByRole('heading', { name: 'Escolhe a tua classe' })).toBeVisible()
  await page.getByRole('button', { name: 'Fighter' }).first().click()
  await page.getByRole('button', { name: 'Próximo →' }).click()

  // Background step
  await expect(page.getByRole('heading', { name: 'Escolhe o teu background' })).toBeVisible()
  await page.getByRole('button', { name: /^Sage/ }).click()
  await page.getByRole('button', { name: 'Próximo →' }).click()

  // Species step
  await expect(page.getByRole('heading', { name: 'Escolhe a tua espécie' })).toBeVisible()
  await page.getByRole('button', { name: /^Human/ }).click()
  await page.getByRole('button', { name: 'Próximo →' }).click()

  // Abilities step
  await expect(page.getByRole('heading', { name: 'Distribui os atributos' })).toBeVisible()
  const selects = page.locator('select')
  const pool = [15, 14, 13, 12, 10, 8]
  const count = await selects.count()
  for (let i = 0; i < count; i += 1) {
    await selects.nth(i).selectOption(String(pool[i]))
  }
  await page.getByRole('button', { name: 'Próximo →' }).click()

  // Equipment step (no gear data in the seed -> trivially valid)
  await expect(page.getByRole('heading', { name: 'Equipamento' })).toBeVisible()
  await page.getByRole('button', { name: 'Próximo →' }).click()

  // Review step + save
  await expect(page.getByRole('heading', { name: 'Revisão final' })).toBeVisible()
  await page.getByPlaceholder(/Ex\.: Thorin Oakenshield/).fill('E2E Hero')
  await page.getByRole('button', { name: '⚔ Forjar personagem' }).click()

  // Character list now shows the saved hero
  await expect(page).toHaveURL(/#\/characters\//)
  await expect(page.getByRole('heading', { name: 'E2E Hero' })).toBeVisible()
})