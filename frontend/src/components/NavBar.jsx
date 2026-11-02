
import { GitBranch } from "lucide-react";

function Navbar() {
  return (
    <nav className="flex items-center justify-between px-8 py-5 bg-white  shadow-md">
      <div className="flex items-center gap-3">
        <GitBranch size={28} className="text-blue-900" />
        <h1 className="text-xl font-bold text-black">
          Dev<span className="text-blue-900">Equity</span>
        </h1>
      </div>

      <div className="flex items-center gap-8 text-blue-900 font-medium">
        <a href="/" className="hover:text-blue-400">Home</a>
        <a href="#features" className="hover:text-blue-400">Features</a>
        <a href="#about" className="hover:text-blue-400">About</a>
      </div>
    </nav>
  );
}

export default Navbar;
