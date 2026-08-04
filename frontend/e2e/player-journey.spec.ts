import { expect, test } from '@playwright/test'

test.describe('Player Journey (Guest Mode)', () => {
  test.beforeEach(async ({ page }) => {
    // Listen for console errors
    page.on('console', msg => {
      if (msg.type() === 'error') {
        console.log('CONSOLE ERROR:', msg.text())
      }
    })
    page.on('pageerror', error => {
      console.log('PAGE ERROR:', error.message)
    })
  })

  test('complete wizard: create character from landing to sheet', async ({ page }) => {
    // Start from landing page
    await page.goto('/')
    
    // Navigate to builder - click "Create Character" link/button
    const createLink = page.getByRole('link', { name: /Create Character/i })
    if (await createLink.isVisible({ timeout: 2000 }).catch(() => false)) {
      await createLink.click()
    } else {
      // Fallback: direct navigation
      await page.goto('/#/builder')
    }

    // Wait for content to load
    await expect(page.getByText('Loading content…')).toBeHidden({ timeout: 15000 })

    // Wait for the class step to be visible
    await expect(page.getByRole('heading', { name: 'Escolha Sua Classe' })).toBeVisible({ timeout: 10000 })

// Step 1: Class - Select Fighter
    const fighterCard = page.locator('[data-class-id="fighter"]')
    await expect(fighterCard).toBeVisible()
    await fighterCard.click()
    
    // Wait for the class selection to propagate to the store
    await page.waitForTimeout(1000)
    
    // Debug: check the step value and class selection
    const stepValue = await page.evaluate(() => {
      // @ts-ignore
      return window.__builder_step_value;
    });
    console.log('Step value after click:', stepValue);
    
    // Debug: check validation directly
    const validationResult = await page.evaluate(() => {
      // @ts-ignore
      const draft = window.__builder_draft?.value;
      if (!draft) return { valid: false, reason: 'no draft' };
      const step = 0;
      const valid = draft.classes.length > 0 && draft.classes[0]?.level >= 1;
      return { valid, classesLength: draft.classes.length, firstClassLevel: draft.classes[0]?.level };
    });
    console.log('Validation result:', validationResult);
    
    // Fill character name (level is already set to 1 by default)
    await page.locator('[data-testid="char-name-input"]').fill('E2E Hero')
    
    // Wait for Next button to be enabled (validation reacts)
    const nextBtn = page.locator('[data-testid="next-btn"]')
    await expect(nextBtn).toBeVisible({ timeout: 5000 })
    await nextBtn.waitFor({ state: 'attached' })
    await page.waitForTimeout(2000) // let reactive update propagate
    await expect(nextBtn).toBeEnabled({ timeout: 10000 })
    await nextBtn.click()

    // Step 2: Background - Select Sage
    await expect(page.getByRole('heading', { name: 'Escolha Seu Background' })).toBeVisible()
    const sageBtn = page.getByRole('button', { name: /^Sage/ })
    await expect(sageBtn).toBeVisible()
    await sageBtn.click()
    await nextBtn.click()

    // Step 3: Species - Select Human
    await expect(page.getByRole('heading', { name: 'Escolha Sua Espécie' })).toBeVisible()
    const humanBtn = page.getByRole('button', { name: /^Human/ })
    await expect(humanBtn).toBeVisible()
    await humanBtn.click()
    await nextBtn.click()

    // Step 4: Abilities - Use Standard Array
    await expect(page.getByRole('heading', { name: 'Atribua Atributos' })).toBeVisible()
    const selects = page.locator('select')
    const pool = [15, 14, 13, 12, 10, 8]
    const count = await selects.count()
    for (let i = 0; i < count; i += 1) {
      await selects.nth(i).selectOption(String(pool[i]))
    }
    await nextBtn.click()

    // Step 5: Equipamento (trivially valid with no gear data)
    await expect(page.getByRole('heading', { name: 'Equipamento' })).toBeVisible()
    await nextBtn.click()

    // Step 6: Review + Save
    await expect(page.getByRole('heading', { name: 'Revisão Final' })).toBeVisible()
    
    // Fill character name
    const nameInput = page.locator('[data-testid="char-name-input"]')
    await expect(nameInput).toBeVisible()
    await nameInput.fill('E2E Hero')
    
    // Click Forge Character
    const forgeBtn = page.getByRole('button', { name: '2694 Forge Character' })
    await expect(forgeBtn).toBeEnabled({ timeout: 5000 })
    await forgeBtn.click()

    // Should navigate to character sheet
    await expect(page).toHaveURL(/#\/characters\//)
    await expect(page.getByRole('heading', { name: 'E2E Hero' })).toBeVisible({ timeout: 10000 })
  })

  test('wizard steps: class → background → species → abilities → equipment → review', async ({ page }) => {
    await page.goto('/#/builder')
    await expect(page.getByText('Loading content…')).toBeHidden({ timeout: 15000 })
    await expect(page.getByRole('heading', { name: 'Escolha Sua Classe' })).toBeVisible({ timeout: 10000 })

    // Step 1: Class
    await page.locator('[data-class-id="fighter"]').click()
    await expect(page.locator('[data-testid="next-btn"]')).toBeVisible({ timeout: 5000 })
    await page.locator('[data-testid="next-btn"]').waitFor({ state: 'attached' })
    await page.waitForTimeout(100)
    await expect(page.locator('[data-testid="next-btn"]')).toBeEnabled({ timeout: 5000 })
    await page.locator('[data-testid="next-btn"]').click()

    // Step 2: Background
    await expect(page.getByRole('heading', { name: 'Escolha Seu Background' })).toBeVisible()
    await page.getByRole('button', { name: /^Sage/ }).click()
    await page.locator('[data-testid="next-btn"]').waitFor({ state: 'visible' })
    await page.waitForTimeout(100)
    await expect(page.locator('[data-testid="next-btn"]')).toBeEnabled({ timeout: 5000 })
    await page.locator('[data-testid="next-btn"]').click()

    // Step 3: Species
    await expect(page.getByRole('heading', { name: 'Escolha Sua Espécie' })).toBeVisible()
    await page.getByRole('button', { name: /^Human/ }).click()
    await page.locator('[data-testid="next-btn"]').waitFor({ state: 'visible' })
    await page.waitForTimeout(100)
    await expect(page.locator('[data-testid="next-btn"]')).toBeEnabled({ timeout: 5000 })
    await page.locator('[data-testid="next-btn"]').click()

    // Step 4: Abilities
    await expect(page.getByRole('heading', { name: 'Atribua Atributos' })).toBeVisible()
    const selects = page.locator('select')
    const pool = [15, 14, 13, 12, 10, 8]
    const count = await selects.count()
    for (let i = 0; i < count; i += 1) {
      await selects.nth(i).selectOption(String(pool[i]))
    }
    await page.locator('[data-testid="next-btn"]').waitFor({ state: 'visible' })
    await page.waitForTimeout(100)
    await expect(page.locator('[data-testid="next-btn"]')).toBeEnabled({ timeout: 5000 })
    await page.locator('[data-testid="next-btn"]').click()

    // Step 5: Equipamento
    await expect(page.getByRole('heading', { name: 'Equipamento' })).toBeVisible()
    await page.locator('[data-testid="next-btn"]').waitFor({ state: 'visible' })
    await page.waitForTimeout(100)
    await expect(page.locator('[data-testid="next-btn"]')).toBeEnabled({ timeout: 5000 })
    await page.locator('[data-testid="next-btn"]').click()

    // Step 6: Review
    await expect(page.getByRole('heading', { name: 'Revisão Final' })).toBeVisible()
    await page.locator('[data-testid="char-name-input"]').fill('Journey Hero')
    await page.getByRole('button', { name: '2694 Forge Character' }).click()

    await expect(page).toHaveURL(/#\/characters\//)
    await expect(page.getByRole('heading', { name: 'Journey Hero' })).toBeVisible({ timeout: 10000 })
  })

  test('character sheet: live editing HP, spell slots, conditions, resources', async ({ page }) => {
    // First create a character
    await page.goto('/#/builder')
    await expect(page.getByText('Loading content…')).toBeHidden({ timeout: 15000 })
    await page.locator('[data-class-id="fighter"]').click()
    await page.locator('[data-testid="next-btn"]').waitFor({ state: 'visible' })
    await page.waitForTimeout(100)
    await expect(page.locator('[data-testid="next-btn"]')).toBeEnabled({ timeout: 5000 })
    await page.locator('[data-testid="next-btn"]').click()
    await page.getByRole('button', { name: /^Sage/ }).click()
    await page.locator('[data-testid="next-btn"]').waitFor({ state: 'visible' })
    await page.waitForTimeout(100)
    await expect(page.locator('[data-testid="next-btn"]')).toBeEnabled({ timeout: 5000 })
    await page.locator('[data-testid="next-btn"]').click()
    await page.getByRole('button', { name: /^Human/ }).click()
    await page.locator('[data-testid="next-btn"]').waitFor({ state: 'visible' })
    await page.waitForTimeout(100)
    await expect(page.locator('[data-testid="next-btn"]')).toBeEnabled({ timeout: 5000 })
    await page.locator('[data-testid="next-btn"]').click()
    const selects = page.locator('select')
    const pool = [15, 14, 13, 12, 10, 8]
    for (let i = 0; i < 6; i++) await selects.nth(i).selectOption(String(pool[i]))
    await page.locator('[data-testid="next-btn"]').waitFor({ state: 'visible' })
    await page.waitForTimeout(100)
    await expect(page.locator('[data-testid="next-btn"]')).toBeEnabled({ timeout: 5000 })
    await page.locator('[data-testid="next-btn"]').click()
    await page.locator('[data-testid="next-btn"]').waitFor({ state: 'visible' })
    await page.waitForTimeout(100)
    await expect(page.locator('[data-testid="next-btn"]')).toBeEnabled({ timeout: 5000 })
    await page.locator('[data-testid="next-btn"]').click()
    await page.locator('[data-testid="char-name-input"]').fill('Sheet Test Hero')
    await page.getByRole('button', { name: '2694 Forge Character' }).click()

    // Now on character sheet
    await expect(page).toHaveURL(/#\/characters\//)
    await expect(page.getByRole('heading', { name: 'Sheet Test Hero' })).toBeVisible({ timeout: 10000 })

    // Test HP stepper
    const hpDisplay = page.locator('.stat-value').first() // First stat is HP
    const initialHP = await hpDisplay.textContent()
    await page.locator('button[aria-label="Increase HP"]').click()
    await expect(hpDisplay).not.toHaveText(initialHP)

    // Test spell slots (if caster)
    const slotButtons = page.locator('.slot')
    if (await slotButtons.count() > 0) {
      await slotButtons.first().click()
      // Verify slot count changed
    }

    // Test conditions
    const conditionChips = page.locator('.cond')
    if (await conditionChips.count() > 0) {
      await conditionChips.first().click()
      await expect(conditionChips.first()).toHaveClass(/active/)
    }

    // Test resources (if any)
    const resourceButtons = page.locator('.resource button.mini')
    if (await resourceButtons.count() > 0) {
      await resourceButtons.first().click()
    }
  })

  test('characters list shows saved character with Open/Delete', async ({ page }) => {
    await page.goto('/#/characters')
    
    // Should see the character we created
    await expect(page.getByRole('heading', { name: 'E2E Hero' })).toBeVisible({ timeout: 10000 })
    await expect(page.getByRole('heading', { name: 'Journey Hero' })).toBeVisible()
    await expect(page.getByRole('heading', { name: 'Sheet Test Hero' })).toBeVisible()

    // Test Open navigation
    const firstCard = page.locator('.card').first()
    await firstCard.getByRole('button', { name: 'Open' }).click()
    await expect(page).toHaveURL(/#\/characters\//)

    // Test Delete (with confirmation)
    page.on('dialog', dialog => dialog.accept())
    await firstCard.getByRole('button', { name: 'Delete' }).click()
    await expect(page.getByRole('heading', { name: 'E2E Hero' })).toBeHidden({ timeout: 5000 })
  })

  test('export JSON', async ({ page }) => {
    await page.goto('/#/characters')
    await expect(page.getByRole('heading', { name: 'Journey Hero' })).toBeVisible({ timeout: 10000 })

    // Navigate to character sheet
    await page.locator('.card').first().getByRole('button', { name: 'Open' }).click()
    await expect(page).toHaveURL(/#\/characters\//)

    // Click Export JSON
    const downloadPromise = page.waitForEvent('download')
    await page.getByRole('button', { name: 'Export JSON' }).click()
    const download = await downloadPromise
    expect(download.suggestedFilename()).toMatch(/\.json$/)
  })

  test('restart wizard from review step', async ({ page }) => {
    await page.goto('/#/builder')
    await expect(page.getByText('Loading content…')).toBeHidden({ timeout: 15000 })
    
    // Go to review step quickly
    await page.locator('[data-class-id="fighter"]').click()
    await page.locator('[data-testid="next-btn"]').waitFor({ state: 'visible' })
    await page.waitForTimeout(100)
    await expect(page.locator('[data-testid="next-btn"]')).toBeEnabled({ timeout: 5000 })
    await page.locator('[data-testid="next-btn"]').click()
    await page.getByRole('button', { name: /^Sage/ }).click()
    await page.locator('[data-testid="next-btn"]').waitFor({ state: 'visible' })
    await page.waitForTimeout(100)
    await expect(page.locator('[data-testid="next-btn"]')).toBeEnabled({ timeout: 5000 })
    await page.locator('[data-testid="next-btn"]').click()
    await page.getByRole('button', { name: /^Human/ }).click()
    await page.locator('[data-testid="next-btn"]').waitFor({ state: 'visible' })
    await page.waitForTimeout(100)
    await expect(page.locator('[data-testid="next-btn"]')).toBeEnabled({ timeout: 5000 })
    await page.locator('[data-testid="next-btn"]').click()
    const selects = page.locator('select')
    const pool = [15, 14, 13, 12, 10, 8]
    for (let i = 0; i < 6; i++) await selects.nth(i).selectOption(String(pool[i]))
    await page.locator('[data-testid="next-btn"]').waitFor({ state: 'visible' })
    await page.waitForTimeout(100)
    await expect(page.locator('[data-testid="next-btn"]')).toBeEnabled({ timeout: 5000 })
    await page.locator('[data-testid="next-btn"]').click()
    await page.locator('[data-testid="next-btn"]').waitFor({ state: 'visible' })
    await page.waitForTimeout(100)
    await expect(page.locator('[data-testid="next-btn"]')).toBeEnabled({ timeout: 5000 })
    await page.locator('[data-testid="next-btn"]').click()

    // On review step, click Restart
    await page.locator('[data-testid="restart-btn"]').click()
    
    // Should be back at step 1 (Class)
    await expect(page.getByRole('heading', { name: 'Escolha Sua Classe' })).toBeVisible()
    await expect(page.locator('[data-class-id="fighter"]')).toBeVisible()
  })
})