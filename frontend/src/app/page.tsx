import RegistrationForm from '@/components/RegistrationForm'

export default function HomePage() {
  const eventName = process.env.NEXT_PUBLIC_EVENT_NAME || 'Tau-Tau Run Fun Run 5K';
  const eventDate = process.env.NEXT_PUBLIC_EVENT_DATE || 'February 15, 2026';
  const eventLocation = process.env.NEXT_PUBLIC_EVENT_LOCATION || 'Gelora Bung Karno Stadium, Jakarta';
  const eventDescription = process.env.NEXT_PUBLIC_EVENT_DESCRIPTION || 'Join us for an exciting 5K fun run event!';

  return (
    <main className="min-h-screen bg-gradient-to-br from-primary-light via-white to-secondary-light">
      {/* Hero Section */}
      <div className="bg-primary text-white py-16">
        <div className="container mx-auto px-4 text-center">
          <h1 className="text-5xl md:text-6xl font-bold mb-4">
            {eventName}
          </h1>
          <p className="text-xl md:text-2xl mb-2">
            📅 {eventDate}
          </p>
          <p className="text-lg md:text-xl">
            📍 {eventLocation}
          </p>
        </div>
      </div>

      {/* Main Content */}
      <div className="container mx-auto px-4 py-12">
        <div className="grid md:grid-cols-2 gap-12 max-w-6xl mx-auto">
          {/* Event Information */}
          <div className="space-y-6">
            <div className="card">
              <h2 className="text-3xl font-bold text-primary mb-4">
                About the Event
              </h2>
              <p className="text-gray-700 mb-4">
                {eventDescription}
              </p>
              <div className="space-y-3">
                <div className="flex items-start">
                  <span className="text-2xl mr-3">🏃</span>
                  <div>
                    <h3 className="font-semibold text-lg">5K Distance</h3>
                    <p className="text-gray-600">Perfect for runners of all levels</p>
                  </div>
                </div>
                <div className="flex items-start">
                  <span className="text-2xl mr-3">🎽</span>
                  <div>
                    <h3 className="font-semibold text-lg">Event Kit</h3>
                    <p className="text-gray-600">Race bib, t-shirt, and medal for all finishers</p>
                  </div>
                </div>
                <div className="flex items-start">
                  <span className="text-2xl mr-3">🏆</span>
                  <div>
                    <h3 className="font-semibold text-lg">Prizes & Awards</h3>
                    <p className="text-gray-600">Awards for top finishers in each category</p>
                  </div>
                </div>
                <div className="flex items-start">
                  <span className="text-2xl mr-3">🎉</span>
                  <div>
                    <h3 className="font-semibold text-lg">Post-Race Celebration</h3>
                    <p className="text-gray-600">Music, food, and fun activities</p>
                  </div>
                </div>
              </div>
            </div>

            <div className="card bg-accent-light">
              <h3 className="text-xl font-bold text-gray-800 mb-3">
                📝 Registration Process
              </h3>
              <ol className="list-decimal list-inside space-y-2 text-gray-700">
                <li>Fill out the registration form</li>
                <li>Receive registration confirmation via email</li>
                <li>Complete payment (details in email)</li>
                <li>Get payment confirmation email</li>
                <li>You're ready to run! 🏃‍♂️</li>
              </ol>
            </div>
          </div>

          {/* Registration Form */}
          <div>
            <div className="card sticky top-6">
              <h2 className="text-3xl font-bold text-secondary mb-6">
                Register Now
              </h2>
              <RegistrationForm />
            </div>
          </div>
        </div>
      </div>

      {/* Footer */}
      <footer className="bg-gray-800 text-white py-8 mt-16">
        <div className="container mx-auto px-4 text-center">
          <p className="text-gray-300">
            © 2025 Tau-Tau Run. All rights reserved.
          </p>
          <p className="text-gray-400 text-sm mt-2">
            Questions? Contact us at admin@tautaurun.com
          </p>
        </div>
      </footer>
    </main>
  )
}
