
function Hero({ children }) {
  return (
    <section className="min-h-screen flex flex-col items-center text-center px-6 pt-20 pb-12 bg-blue-200">
      <p className="text-blue-900 font-extrabold tracking-widest text-6xl mb-5">
        DevEquity
      </p>

      <h1 className="text-3xl  font-bold text-black max-w-4xl leading-tight">
        Understand Your Team's
        <span className="block text-blue-900">
          Contributions
        </span>
      </h1>

      <p className="text-slate-900 text-lg max-w-2xl mt-6 mb-10 leading-relaxed">
        Analyze repository activity, track developer contributions,
        and discover meaningful insights into your team's work.
      </p>

      {children}
    </section>
  );
}

export default Hero;
