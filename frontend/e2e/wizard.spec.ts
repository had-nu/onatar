import { expect, test } from '@playwright/test'

test('landing page shows the hero and navigation', async ({ page }) => {
  await page.goto('/')
  await expect(page.getByRole('heading', { level: 1 })).toHaveText('Forge Your Hero')
  await expect(page.getByRole('navigation', { name: 'Main navigation' })).toBeVisible()
  await expect(page.getByRole('link', { name: 'Content' })).toBeVisible()
})

test('wizard flow: create a character end to end', async ({ page }) => {
  const errors: string[] = []
  page.on('console', msg => {
    if (msg.type() === 'error') {
      console.log('CONSOLE ERROR:', msg.text())
    }
  })
  page.on('pageerror', error => {
    console.log('PAGE ERROR:', error.message)
  })

  await page.goto('/#/builder')

  // Wait for content to load (give it up to 15 seconds)
  await expect(page.getByText('Loading content…')).toBeHidden({ timeout: 15000 })

  // Wait for the class step to be visible
  await expect(page.getByRole('heading', { name: 'Choose Your Class' })).toBeVisible({ timeout: 10000 })

  // Debug: log the page content
  const content = await page.content()
  console.log('Page content length:', content.length)
  console.log('Has Choose Your Class:', content.includes('Choose Your Class'))
  console.log('Has Loading:', content.includes('Loading content'))
  console.log('Has Fighter:', content.includes('Fighter'))
  console.log('Has class-grid:', content.includes('class-grid'))
  console.log('Has step-class:', content.includes('step-class'))
  console.log('Has builder-page:', content.includes('builder-page'))
  console.log('Has builder-content:', content.includes('builder-content'))
  console.log('Has steps:', content.includes('StepsBar'))

  // Check if classes array is loaded
  const classesCount = await page.evaluate(() => {
    // @ts-ignore
    return window.__builder_classes_length || 0
  })
  console.log('Classes count:', classesCount)

  // Check what step we're on
  const stepValue = await page.evaluate(() => {
    // @ts-ignore
    return (window as any).__builder_step_value || -1
  })
  console.log('Step value:', stepValue)

  // Check what component is rendered
  const stepComponent = await page.evaluate(() => {
    // @ts-ignore
    return (window as any).__builder_step_component || 'unknown'
  })
  console.log('Step component:', stepComponent)

  // Check what headings exist
  const headings = await page.evaluate(() => {
    const headings = document.querySelectorAll('h1, h2, h3, h4, h5, h6')
    return Array.from(headings).map(h => ({ text: h.textContent?.trim(), tag: h.tagName, class: h.className, visible: h.offsetParent !== null }))
  })
  console.log('Headings:', JSON.stringify(headings, null, 2))

  // Check if step-class div exists and its children
  const stepClassChildren = await page.evaluate(() => {
    const stepClass = document.querySelector('.step-class')
    if (!stepClass) return null
    return {
      html: stepClass.innerHTML.substring(0, 500),
      children: Array.from(stepClass.children).map(c => ({ tag: c.tagName, class: c.className, text: c.textContent?.trim().substring(0, 50) }))
    }
  })
  console.log('Step class children:', JSON.stringify(stepClassChildren, null, 2))

  // Class step
  await expect(page.getByRole('heading', { name: 'Choose Your Class' })).toBeVisible()
  await page.locator('.class-card:has-text("Fighter")').first().click()
  
  // Wait for Next button to be enabled
  await expect(page.getByRole('button', { name: 'Next →' })).toBeEnabled({ timeout: 5000 })
  await page.getByRole('button', { name: 'Next →' }).click()

  // Background step
  await expect(page.getByRole('heading', { name: 'Choose Your Background' })).toBeVisible()
  await page.getByRole('button', { name: /^Sage/ }).click()
  await page.getByRole('button', { name: 'Next →' }).click()

  // Species step
  await expect(page.getByRole('heading', { name: 'Choose Your Species' })).toBeVisible()
  await page.getByRole('button', { name: /^Human/ }).click()
  await page.getByRole('button', { name: 'Next →' }).click()

  // Abilities step
  await expect(page.getByRole('heading', { name: 'Assign Abilities' })).toBeVisible()
  const selects = page.locator('select')
  const pool = [15, 14, 13, 12, 10, 8]
  const count = await selects.count()
  for (let i = 0; i < count; i += 1) {
    await selects.nth(i).selectOption(String(pool[i]))
  }
  await page.getByRole('button', { name: 'Next →' }).click()

  // Equipment step (no gear data in the seed -> trivially valid)
  await expect(page.getByRole('heading', { name: 'Equipment' })).toBeVisible()
  await page.getByRole('button', { name: 'Next →' }).click()

  // Review step + save
  await expect(page.getByRole('heading', { name: 'Final Review' })).toBeVisible()
  await page.getByPlaceholder(/e\.g\., Thorin Oakenshield/).fill('E2E Hero')
  await page.getByRole('button', { name: '⚔ Forge Character' }).click()

  // Character list now shows the saved hero
  await expect(page).toHaveURL(/#\/characters\//)
  await expect(page.getByRole('heading', { name: 'E2E Hero' })).toBeVisible()
})