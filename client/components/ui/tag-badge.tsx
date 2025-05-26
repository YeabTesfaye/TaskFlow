import { Tag } from "@/types/task";
import { cn } from "@/lib/utils";

interface TagBadgeProps {
  tag: Tag;
  className?: string;
  size?: "sm" | "default";
}

export function TagBadge({ tag, className, size = "default" }: TagBadgeProps) {
  // Calculate text color based on background color luminance
  const getContrastTextColor = (hexColor: string) => {
    // Convert hex to RGB
    const r = parseInt(hexColor.slice(1, 3), 16);
    const g = parseInt(hexColor.slice(3, 5), 16);
    const b = parseInt(hexColor.slice(5, 7), 16);
    
    // Calculate luminance
    const luminance = (0.299 * r + 0.587 * g + 0.114 * b) / 255;
    
    // Return black for bright colors, white for dark colors
    return luminance > 0.5 ? "#000000" : "#ffffff";
  };

  const style = {
    backgroundColor: `${tag.color}20`, // 20% opacity
    color: tag.color,
    borderColor: `${tag.color}40`, // 40% opacity
  };

  return (
    <span
      className={cn(
        "inline-block rounded-full border text-xs font-medium",
        size === "sm" ? "px-2 py-0.5 text-[10px]" : "px-2.5 py-0.5",
        className
      )}
      style={style}
    >
      {tag.name}
    </span>
  );
}