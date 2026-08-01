import { expect, test } from '@playwright/test'

test('landing page shows the hero and navigation', async ({ page }) => {
  await page.goto('/')
  await expect(page.getByRole('heading', { level: 1 })).toHaveText('Onatar')
  await expect(page.getByRole('navigation', { name: 'Navegação principal' })).toBeVisible()
  await expect(page.getByRole('link', { name: 'Personagens', exact: true })).toBeVisible()
})

test('wizard flow: create a character end to end', async ({ page }) => {
  await page.goto('/#/builder')

  // Class step
  await expect(page.getByRole('heading', { name: 'Escolhe as classes' })).toBeVisible()
  await page.getByRole('button', { name: 'Escolher' }).first().click()
  await page.getByRole('button', { name: 'Continuar →' }).click()

  // Background step
  await expect(page.getByRole('heading', { name: 'Escolhe o background' })).toBeVisible()
  await page.getByRole('button', { name: /^Sage/ }).click()
  await page.getByRole('button', { name: 'Continuar →' }).click()

  // Species step
  await expect(page.getByRole('heading', { name: 'Escolhe a espécie' })).toBeVisible()
  await page.getByRole('button', { name: /^Human/ }).click()
  await page.getByRole('button', { name: 'Continuar →' }).click()

  // Abilities step
  await expect(page.getByRole('heading', { name: 'Atributos' })).toBeVisible()
  const selects = page.locator('select')
  const pool = [15, 14, 13, 12, 10, 8]
  const count = await selects.count()
  for (let i = 0; i < count; i += 1) {
    await selects.nth(i).selectOption(String(pool[i]))
  }
  await page.getByRole('button', { name: 'Continuar →' }).click()

  // Equipment step (no gear data in the seed -> trivially valid)
  await expect(page.getByRole('heading', { name: 'Equipamento' })).toBeVisible()
  await page.getByRole('button', { name: 'Continuar →' }).click()

  // Review step + save
  await expect(page.getByRole('heading', { name: 'Revisão' })).toBeVisible()
  await page.getByPlaceholder(/Ex\.: Onatar/).fill('E2E Hero')
  await page.getByRole('button', { name: 'Guardar personagem' }).click()

  // Character list now shows the saved hero
  await expect(page).toHaveURL(/#\/characters\//)
  await expect(page.getByRole('heading', { name: 'E2E Hero' })).toBeVisible()
})
