const pages = Array.from(document.querySelectorAll('.page'))
const navItems = Array.from(document.querySelectorAll('.nav-item'))
navItems.forEach(btn => btn.addEventListener('click', () => {
  navItems.forEach(b => b.classList.remove('active'))
  btn.classList.add('active')
  const target = btn.dataset.target
  pages.forEach(p => p.classList.toggle('hidden', p.id !== target))
}))

// Modal (Rule Wizard)
const wizard = document.getElementById('ruleWizard')
document.getElementById('createRuleBtn').addEventListener('click', () => {
  wizard.classList.remove('hidden')
})
document.getElementById('closeWizard').addEventListener('click', () => {
  wizard.classList.add('hidden')
})
document.getElementById('prevStep').addEventListener('click', () => {
  step(-1)
})
document.getElementById('nextStep').addEventListener('click', () => {
  step(1)
})

let currentStep = 0
function step(delta) {
  const steps = Array.from(document.querySelectorAll('.steps .step'))
  currentStep = Math.min(Math.max(0, currentStep + delta), steps.length - 1)
  steps.forEach((s, i) => s.classList.toggle('active', i === currentStep))
}

// Mini charts in KPI cards
Array.from(document.querySelectorAll('.mini-chart')).forEach(c => drawSparkline(c))

function drawSparkline(canvas) {
  const ctx = canvas.getContext('2d')
  const w = canvas.width, h = canvas.height
  ctx.clearRect(0,0,w,h)
  const gradient = ctx.createLinearGradient(0,0,w,0)
  gradient.addColorStop(0, '#00e5ff')
  gradient.addColorStop(1, '#7c4dff')
  ctx.strokeStyle = gradient
  ctx.lineWidth = 2
  const points = Array.from({length: 24}, (_, i) => {
    const x = (i/(24-1))*w
    const y = h - (Math.sin(i*0.35)+1)/2 * (h*0.8) - h*0.1
    return {x,y}
  })
  ctx.beginPath()
  points.forEach((p, i) => i ? ctx.lineTo(p.x, p.y) : ctx.moveTo(p.x, p.y))
  ctx.stroke()
}

// Quality trend chart
const qc = document.getElementById('qualityChart')
if (qc) {
  const ctx = qc.getContext('2d')
  const w = qc.width, h = qc.height
  const gridColor = 'rgba(156,195,255,0.15)'
  ctx.clearRect(0,0,w,h)
  // grid
  for (let x=0; x<=w; x+=72) {
    ctx.strokeStyle = gridColor
    ctx.beginPath(); ctx.moveTo(x,0); ctx.lineTo(x,h); ctx.stroke()
  }
  for (let y=0; y<=h; y+=52) {
    ctx.strokeStyle = gridColor
    ctx.beginPath(); ctx.moveTo(0,y); ctx.lineTo(w,y); ctx.stroke()
  }
  // line
  const grad = ctx.createLinearGradient(0,0,w,0)
  grad.addColorStop(0,'#00e5ff')
  grad.addColorStop(1,'#7c4dff')
  ctx.strokeStyle = grad
  ctx.lineWidth = 2
  const pts = Array.from({length: 36}, (_, i) => {
    const x = (i/(36-1))*w
    const y = h - (Math.sin(i*0.25)+1)/2 * (h*0.7) - h*0.15
    return {x,y}
  })
  ctx.beginPath()
  pts.forEach((p, i) => i ? ctx.lineTo(p.x, p.y) : ctx.moveTo(p.x, p.y))
  ctx.stroke()
  // glow
  ctx.save()
  ctx.shadowColor = 'rgba(124,77,255,0.35)'
  ctx.shadowBlur = 24
  ctx.beginPath()
  pts.forEach((p, i) => i ? ctx.lineTo(p.x, p.y) : ctx.moveTo(p.x, p.y))
  ctx.stroke()
  ctx.restore()
}

// Theme toggle (visual accent only)
document.getElementById('theme').addEventListener('change', (e) => {
  const on = e.target.checked
  document.body.style.background = on
    ? 'radial-gradient(1200px 600px at 0% 0%, rgba(0,229,255,0.22), transparent 60%),radial-gradient(1000px 600px at 100% 100%, rgba(124,77,255,0.22), transparent 60%), #07090d'
    : ''
})

