import { Link } from 'react-router-dom'
import heroImage from '../assets/landing-hero.png'
import './LandingPage.css'

function LandingPage() {
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
          <h1 className="landing-title">
            Ambiguity
            <br />
            Detector
          </h1>

          <div className="landing-buttons">
            <Link className="landing-button landing-button--primary" to="/login">
              Sign in
            </Link>

            <Link className="landing-button landing-button--secondary" to="/register">
              Sign up
            </Link>
          </div>
        </section>

        <section className="landing-image-wrapper" aria-hidden="true">
          <img className="landing-image" src={heroImage} alt="Ambiguity Detector preview" />
        </section>
      </main>
    </div>
  )
}

export default LandingPage