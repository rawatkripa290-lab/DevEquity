
import { GitBranch } from "lucide-react";

function Navbar() {
  return (
    <nav className="sticky top-0 z-50 flex items-center justify-between border-b border-slate-200 bg-white/95 px-6 py-4 shadow-sm backdrop-blur-md md:px-10">

      {/* Logo */}
      <a href="#home" className="flex items-center gap-3">
        <div className="rounded-lg bg-blue-100 p-2">
          <GitBranch size={25} className="text-blue-700" />
        </div>

        <h1 className="text-2xl font-extrabold tracking-tight text-slate-900">
          Dev<span className="text-blue-700">Equity</span>
        </h1>
      </a>

      {/* Navigation Links */}
      <div className="flex items-center gap-5 text-medium font-semibold text-slate-600 md:gap-9">
        <a href="#home" className="transition-colors hover:text-blue-700">
          Home
        </a>

        <a href="#features" className="transition-colors hover:text-blue-700">
          Features
        </a>

        <a href="#about" className="transition-colors hover:text-blue-700">
          About
        </a>

        <a
          href="#home"
          className="text-xl rounded-lg bg-blue-700 px-5 py-2.5 text-white hover:bg-blue-800"
        >
          Analyze Repository
        </a>
      </div>
    </nav>
  );
}

export default Navbar;
