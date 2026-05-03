import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import heroImage from '../assets/hero.png'
import './LandingPage.css'

const sloganWords = ['лучше', 'понятнее', 'точнее', 'красивее', 'аккуратнее']

function LandingPage() {
  const [wordIndex, setWordIndex] = useState(0)

  useEffect(() => {
    const intervalId = setInterval(() => {
      setWordIndex((prev) => (prev + 1) % sloganWords.length)
    }, 2000)

    return () => clearInterval(intervalId)
  }, [])

  return (
    <div className="landing-page">
      <header className="landing-header">
        <div className="landing-header-actions">
          <Link className="landing-header-link" to="/login">
            Sign in
          </Link>

          <Link className="landing-header-link" to="/register">
            Sign up
          </Link>

          <div className="landing-header-avatar" aria-hidden="true" />
        </div>
      </header>

      <main className="landing-main">
        <section className="landing-content">
          <div className="landing-title-block">
            <h1 className="landing-title">
              <span>Ambiguity</span>
              <span>Detector</span>
            </h1>

            <p className="landing-slogan">
              <span className="landing-slogan-static">Сделай ТЗ</span>{' '}
              <span key={sloganWords[wordIndex]} className="landing-slogan-word">
                {sloganWords[wordIndex]}
              </span>
            </p>
          </div>
        </section>

        <section className="landing-image-wrapper" aria-hidden="true">
          <img
            className="landing-image"
            src={heroImage}
            alt="Ambiguity Detector preview"
          />
        </section>
      </main>
    </div>
  )
}

export default LandingPage