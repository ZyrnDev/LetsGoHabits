import { FC } from "react";
import Head from "next/head";
import { Center, CenterProps } from "@chakra-ui/react";
import { Footer, content_min_height, footer_height, padding } from "@components/layouts/main"
import Meta, { MetaOptions } from "@components/meta";
import { motion } from "framer-motion";
// import { basePath } from "@/lib/env";

const transitionYDistance = 25;

const max_content_view_height = `${content_min_height} + 2 *  ${padding} + ${footer_height}`;
const max_height_during_transition = `${max_content_view_height} - ${transitionYDistance}px`;

const variants = {
  hidden: { opacity: 0, y: transitionYDistance, maxHeight: `calc(${max_height_during_transition})` },
  enter: { opacity: 1, y: 0, transitionEnd: { maxHeight: 'none' } },
  exit: { opacity: 0, y: transitionYDistance, maxHeight: `calc(${max_height_during_transition})` },
};

const Animator: FC = ({ children }) => {
  return (
    <motion.div
      initial="hidden"
      animate="enter"
      exit="exit"
      variants={variants}
      style={{ overflowY: "hidden" }}
      transition={{ duration: 0.4, ease: "easeInOut" }}
    >
      {children}
    </motion.div>
  );
};


const Layout: FC<CenterProps & { meta: MetaOptions }> = ({ meta = {}, children, ...props }) => {
  const maxH = props.maxH || props.maxHeight;
  const minH = maxH ? `min(${content_min_height}, ${maxH})` : `calc(${content_min_height})`;
  return (
    <>
      <Head>
        <link rel="icon" href={`/favicon.ico`} />
      </Head>
      <Meta {...meta} />
      <Animator>
        <Center minH={minH} flexDir="column" as="section" flex="1" textAlign="center" {...props}>
          {children}
        </Center>
        <Footer p={padding} h={footer_height}/> 
      </Animator>
    </>
  );
};

export default Layout;