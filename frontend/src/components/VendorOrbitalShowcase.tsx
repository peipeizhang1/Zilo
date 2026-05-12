import { motion } from "framer-motion";

type Vendor = {
  name: string;
  short: string;
  angle: number;
  radius: number;
  color: string;
};

const vendors: Vendor[] = [
  { name: "OpenAI", short: "OA", angle: 0, radius: 130, color: "#4FD1FF" },
  { name: "Anthropic", short: "AN", angle: 36, radius: 130, color: "#79E2A0" },
  { name: "Gemini", short: "GM", angle: 72, radius: 130, color: "#8B7CFF" },
  { name: "通义", short: "TY", angle: 108, radius: 130, color: "#62C3FF" },
  { name: "豆包", short: "DB", angle: 144, radius: 130, color: "#8CF3FF" },
  { name: "DeepSeek", short: "DS", angle: 180, radius: 130, color: "#A6B5FF" },
  { name: "Moonshot", short: "MS", angle: 216, radius: 130, color: "#65D7FF" },
  { name: "Groq", short: "GQ", angle: 252, radius: 130, color: "#7EFFC7" },
  { name: "Azure", short: "AZ", angle: 288, radius: 130, color: "#A8A8FF" },
  { name: "SiliconFlow", short: "SF", angle: 324, radius: 130, color: "#69D9FF" }
];

export function VendorOrbitalShowcase() {
  return (
    <div className="orbital-root">
      <motion.div
        className="orbital-ring"
        animate={{ rotate: 360 }}
        transition={{ duration: 45, repeat: Infinity, ease: "linear" }}
      >
        {vendors.map((item) => {
          const rad = (item.angle * Math.PI) / 180;
          const x = Math.cos(rad) * item.radius;
          const y = Math.sin(rad) * item.radius;
          return (
            <motion.button
              key={item.name}
              className="vendor-node"
              style={{ transform: `translate(${x}px, ${y}px)`, borderColor: item.color }}
              whileHover={{ scale: 1.12, y: -8 }}
              whileTap={{ scale: 0.96 }}
              title={item.name}
            >
              <span style={{ color: item.color }}>{item.short}</span>
            </motion.button>
          );
        })}
      </motion.div>
      <div className="orbital-core">
        <div className="core-title">Zilo</div>
        <div className="core-sub">Workflow Engine</div>
      </div>
    </div>
  );
}

