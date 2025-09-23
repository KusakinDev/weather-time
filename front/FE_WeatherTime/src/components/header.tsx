"use client";


const Header: React.FC = () => {
 
  return (
    <header className="bg-bgHeaderFooter text-txHeaderFooter py-4 sticky top-0 left-0 w-full z-10 shadow-md">
      <div className="max-w-7xl mx-auto flex justify-between items-center px-6">
        <div className="flex items-center space-x-3">
            <img src="/logo.png" alt="Logo" className="h-8 w-8" />
            <h1 className="text-xl font-semibold">Alexander&apos;s weather-time city</h1>
        </div>
      </div>
    </header>
  );
};

export default Header;
